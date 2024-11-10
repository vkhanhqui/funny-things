package service

import (
	"cloud-gaming/game"
	"encoding/json"
	"log"
	"time"

	"github.com/pion/webrtc/v4"
	"github.com/pion/webrtc/v4/pkg/media"
)

func NewWorker() *Worker {
	return &Worker{}
}

type Worker struct {
}

func (w *Worker) Run() {
	go func() {
		for p := range peerConnCh {
			peerConns = append(peerConns, p)
			log.Println("Added to connections")

			w.onDataChannel(p)
		}
	}()
}

func (w *Worker) onDataChannel(p *PeerConnState) {
	closeSignal := make(chan bool)
	cmdCh := make(chan string)
	gameStateCh := make(chan *game.Snake, 1)
	pixelCh := make(chan []byte)
	encodedFrameCh := make(chan []byte)
	senders := p.PeerConnection().GetSenders()

	pc := p.PeerConnection()
	pc.OnDataChannel(func(dataCh *webrtc.DataChannel) {
		gameLoop := game.NewSnakeLoop(&game.SnakeLoopInit{CommandChannel: cmdCh, SnakeChannel: gameStateCh, CloseSignal: closeSignal})
		go gameLoop.Start()
		go game.StartFrameRenderer(gameStateCh, pixelCh)

		go w.startEncoder(pixelCh, encodedFrameCh)
		go w.startStreaming(encodedFrameCh, senders)
		go w.closeConnection(closeSignal, dataCh, p, gameStateCh, cmdCh, pixelCh, encodedFrameCh)

		w.onMessage(dataCh, cmdCh)
		w.onError(dataCh, closeSignal)
	})
}

func (w *Worker) onError(dataCh *webrtc.DataChannel, closeSignal chan bool) {
	dataCh.OnError(func(err error) {
		if err != nil {
			log.Println(err)
			closeSignal <- true
		}
	})
}

func (w *Worker) closeConnection(closeSignal chan bool, dataCh *webrtc.DataChannel,
	peerConn *PeerConnState, gameStateCh chan *game.Snake, cmdCh chan string,
	pixelCh, encodedFrameCh chan []byte) {
	<-closeSignal

	err := dataCh.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = peerConn.Close()
	if err != nil {
		log.Fatal(err)
	}

	close(gameStateCh)
	close(cmdCh)
	close(pixelCh)
	close(encodedFrameCh)
}

func (w *Worker) onMessage(dataCh *webrtc.DataChannel, cmdCh chan string) {
	dataCh.OnMessage(func(msg webrtc.DataChannelMessage) {
		var message Message
		err := json.Unmarshal(msg.Data, &message)
		if err != nil || message.Type != COMMAND {
			log.Fatal(err)
		}

		cmdCh <- message.Value
	})
}

func (w *Worker) startStreaming(encodedFrameCh chan []byte, senders []*webrtc.RTPSender) {
	for {
		encodedFrame, ok := <-encodedFrameCh
		if !ok {
			break
		}

		for _, s := range senders {
			track := s.Track().(*webrtc.TrackLocalStaticSample)
			err := track.WriteSample(media.Sample{Data: encodedFrame, Duration: time.Second / game.FPS, Timestamp: time.Now()})
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
