package service

import (
	"cloud-gaming/game"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/pion/webrtc/v4/pkg/media/h264reader"
)

const ffmpegCommand = "ffmpeg -hide_banner -loglevel error -f rawvideo -pixel_format rgb24 -video_size %dx%d -framerate %d -i pipe:0 -c:v libx264 -preset ultrafast -tune zerolatency -f h264 pipe:1"

func (w *Worker) startEncoder(canvasCh chan *game.Canvas, encodedFrameCh chan *Streamable) {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("run time panic: %v", x)
		}
	}()

	for {
		start := time.Now()
		rawRGBDataFrame, ok := <-canvasCh
		if !ok {
			break
		}

		cmd := exec.Command("bash", "-c", fmt.Sprintf(ffmpegCommand, game.FRAME_WIDTH, game.FRAME_HEIGHT, game.FPS))
		cmd.Stderr = os.Stderr

		inPipe, err := cmd.StdinPipe()
		if err != nil {
			log.Fatal(err)
		}

		outPipe, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}

		err = cmd.Start()
		if err != nil {
			log.Fatal(err)
		}

		_, err = inPipe.Write(rawRGBDataFrame.Data)
		if err != nil {
			log.Fatal(err)
		}

		inPipe.Close()

		encodedData, err := readH264NALUnits(outPipe)
		if err != nil {
			log.Fatal(err)
		}

		err = cmd.Wait()
		if err != nil {
			log.Fatal(err)
		}

		encodedFrameCh <- &Streamable{Data: encodedData, Timestamp: time.Now()}
		log.Println("startEncoder", time.Since(start))
	}
}

func readH264NALUnits(outPipe io.Reader) ([]byte, error) {
	h264, err := h264reader.NewReader(outPipe)
	if err != nil {
		return nil, fmt.Errorf("failed to create H.264 reader: %v", err)
	}

	var data []byte
	var spsAndPpsCache []byte

	for {
		nal, h264Err := h264.NextNAL()
		if h264Err == io.EOF {
			break
		} else if h264Err != nil {
			return nil, fmt.Errorf("error reading H.264 NAL: %v", h264Err)
		}

		nal.Data = append([]byte{0x00, 0x00, 0x00, 0x01}, nal.Data...)
		if nal.UnitType == h264reader.NalUnitTypeSPS || nal.UnitType == h264reader.NalUnitTypePPS {
			spsAndPpsCache = append(spsAndPpsCache, nal.Data...)
			continue
		} else if nal.UnitType == h264reader.NalUnitTypeCodedSliceIdr {
			nal.Data = append(spsAndPpsCache, nal.Data...)
			spsAndPpsCache = []byte{}
		}

		data = append(data, nal.Data...)
	}

	return data, nil
}
