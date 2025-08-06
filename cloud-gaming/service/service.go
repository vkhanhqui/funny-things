package service

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v4"
)

var (
	peerConns  []*PeerConnState
	peerConnCh = make(chan *PeerConnState)
)

func NewService() *Service {
	return &Service{}
}

type Service struct{}

func (s *Service) HandleMessage(p *PeerConnState) {
	p.PeerConnection().OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		log.Println("Connection state changed:", state)
		if state == webrtc.PeerConnectionStateConnected {
			peerConnCh <- p
			log.Println("Added new peer connection to channel")
		}
	})

	for {
		msg, err := p.ReadMessage()
		if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseAbnormalClosure, websocket.CloseGoingAway) {
			log.Printf("Websocket closed: %v", err)
			break
		} else if err != nil {
			log.Printf("Error reading from websocket: %v", err)
			continue
		}

		switch msg.Type {
		case OFFER:
			offer := msg.Value
			err = p.CreateOffer(offer)
			if err != nil {
				log.Printf("Error when creating offer: %v", err)
				continue
			}

		case ICE:
			ice := msg.Value
			err = p.AddICECandidate(ice)
			if err != nil {
				log.Printf("Error when adding ICE: %v", err)
				continue
			}

		default:
			log.Printf("Unknown message type: %+v", msg)
			continue
		}
	}
}
