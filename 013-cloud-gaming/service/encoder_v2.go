package service

import (
	"cloud-gaming/game"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
	"sync/atomic"
	"time"
)

type Encoder struct {
	encodedFrameCh chan *Streamable
	canvasCh       <-chan *game.Canvas
	closeSignal    <-chan bool
	windowWidth    int
	windowHeight   int
	closed         int32
	wg             sync.WaitGroup
}

type EncoderOptions struct {
	EncodedFrameChannel chan *Streamable
	CanvasChannel       <-chan *game.Canvas
	CloseSignal         <-chan bool
	WindowWidth         int
	WindowHeight        int
}

type Streamable struct {
	Data      []byte
	Timestamp time.Time
}

func NewEncoder(options *EncoderOptions) *Encoder {
	return &Encoder{
		encodedFrameCh: options.EncodedFrameChannel,
		canvasCh:       options.CanvasChannel,
		closeSignal:    options.CloseSignal,
		windowWidth:    options.WindowWidth,
		windowHeight:   options.WindowHeight,
	}
}

func (e *Encoder) isClosed() bool {
	return atomic.LoadInt32(&e.closed) == 1
}

func (e *Encoder) markAsClosed() {
	atomic.StoreInt32(&e.closed, 1)
}

const ffmpegBaseCommand = "ffmpeg -hide_banner -loglevel error -threads 0 -re -f rawvideo -pixel_format rgb24 -video_size %dx%d -framerate %v -r %v -i pipe:0 -pix_fmt yuv420p -c:v h264_videotoolbox -f h264 pipe:1"

func (e *Encoder) Start() {
	ffmpegCommand := fmt.Sprintf(ffmpegBaseCommand, e.windowWidth, e.windowHeight, game.FPS, game.FPS)
	cmd := exec.Command("bash", "-c", ffmpegCommand)
	cmd.Stderr = os.Stderr
	inPipe, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}

	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	e.wg.Add(2)

	go e.writeToFFmpeg(inPipe)
	go e.streamToWebRTCTrack(outPipe)

	go func() {
		_, ok := <-e.closeSignal
		if !ok {
			e.markAsClosed()
			fmt.Println("Closing encoder")

			cmd.Wait()
			e.wg.Wait()
			close(e.encodedFrameCh)
		}
	}()

}

func (e *Encoder) writeToFFmpeg(inPipe io.WriteCloser) {
	defer e.wg.Done()

	for canvas := range e.canvasCh {
		if e.isClosed() {
			return
		}
		select {
		case <-e.encodedFrameCh: // Check if there's a backlog.
			continue
		default:
			_, err := inPipe.Write(canvas.Data)
			if err != nil {
				panic(err)
			}
		}
	}

	inPipe.Close()
}

func (e *Encoder) streamToWebRTCTrack(outPipe io.Reader) {
	defer e.wg.Done()

	buf := make([]byte, 1024*8)
	for {
		if e.isClosed() {
			return
		}

		timestamp := time.Now()
		n, err := outPipe.Read(buf)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println(err)
			continue
		}

		e.encodedFrameCh <- &Streamable{Data: buf[:n], Timestamp: timestamp}
	}
}
