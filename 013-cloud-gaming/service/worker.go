package service

import (
	"cloud-gaming/game"
	"encoding/json"
	"log"

	"github.com/pion/webrtc/v4"
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

	pc := p.PeerConnection()
	pc.OnDataChannel(func(dataCh *webrtc.DataChannel) {
		gameLoop := game.NewSnakeLoop(&game.SnakeLoopInit{CommandChannel: cmdCh, SnakeChannel: gameStateCh, CloseSignal: closeSignal})
		go gameLoop.Start()

		go w.closeConnection(dataCh, p, gameStateCh, cmdCh, closeSignal)

		w.onMessage(dataCh, cmdCh)
	})
}

func (w *Worker) closeConnection(dataCh *webrtc.DataChannel, peerConn *PeerConnState,
	gameStateCh chan *game.Snake, cmdCh chan string, closeSignal chan bool) {
	<-closeSignal
	log.Println("Closing peer connection")

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
}

func (w *Worker) onMessage(dataCh *webrtc.DataChannel, cmdCh chan string) {
	dataCh.OnMessage(func(msg webrtc.DataChannelMessage) {
		var message Message
		err := json.Unmarshal(msg.Data, &message)
		if err != nil || message.Type != COMMAND {
			log.Fatal(err)
		}

		log.Println("Received: ", message.Value)
	})
}
