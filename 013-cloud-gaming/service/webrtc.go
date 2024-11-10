package service

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v4"
)

type MessageType string

const (
	OFFER  MessageType = "offer"
	ICE    MessageType = "ice"
	ANSWER MessageType = "answer"
	ERROR  MessageType = "error"
)

type Message struct {
	Type  MessageType `json:"type"`
	Value any         `json:"value"`
}

type PeerConnState struct {
	peerConnection *webrtc.PeerConnection
	ws             *threadSafeWriter
}

func NewPeerConnection(peerConnection *webrtc.PeerConnection, ws *threadSafeWriter) *PeerConnState {
	return &PeerConnState{peerConnection, ws}
}

func (p *PeerConnState) PeerConnection() *webrtc.PeerConnection {
	return p.peerConnection
}

func (p *PeerConnState) WS() *threadSafeWriter {
	return p.ws
}

func (p *PeerConnState) ReadMessage() (Message, error) {
	var msg Message
	err := p.ws.ReadJSON(&msg)
	return msg, err
}

func (p *PeerConnState) CreateOffer(offer string) error {
	err := p.PeerConnection().SetRemoteDescription(webrtc.SessionDescription{Type: webrtc.SDPTypeOffer, SDP: offer})
	if err != nil {
		return err
	}

	answer, err := p.PeerConnection().CreateAnswer(nil)
	if err != nil {
		return err
	}

	err = p.PeerConnection().SetLocalDescription(answer)
	if err != nil {
		return err
	}

	err = p.WS().WriteJSON(ANSWER, answer.SDP)
	return err
}

func (p *PeerConnState) AddICECandidate(ice string) error {
	return p.peerConnection.AddICECandidate(webrtc.ICECandidateInit{Candidate: ice})
}

type threadSafeWriter struct {
	*websocket.Conn
	sync.Mutex
}

func NewSafeWriter(unsafeConn *websocket.Conn) *threadSafeWriter {
	return &threadSafeWriter{unsafeConn, sync.Mutex{}}
}

func (t *threadSafeWriter) WriteJSON(mt MessageType, value string) error {
	t.Lock()
	defer t.Unlock()
	return t.Conn.WriteJSON(Message{mt, value})
}
