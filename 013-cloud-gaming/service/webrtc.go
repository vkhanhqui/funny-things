package service

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v4"
)

type MessageType string

const (
	OFFER   MessageType = "offer"
	ICE     MessageType = "ice"
	ANSWER  MessageType = "answer"
	COMMAND MessageType = "command"
	ERROR   MessageType = "error"
)

type Message struct {
	Type  MessageType `json:"type"`
	Value string      `json:"value"`
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

func (p *PeerConnState) Close() error {
	return p.peerConnection.Close()
}

func (p *PeerConnState) ReadMessage() (Message, error) {
	var msg Message
	err := p.ws.ReadJSON(&msg)
	return msg, err
}

func (p *PeerConnState) CreateOffer(offer string) error {
	err := p.peerConnection.SetRemoteDescription(webrtc.SessionDescription{Type: webrtc.SDPTypeOffer, SDP: offer})
	if err != nil {
		return err
	}

	answer, err := p.peerConnection.CreateAnswer(nil)
	if err != nil {
		return err
	}

	err = p.peerConnection.SetLocalDescription(answer)
	if err != nil {
		return err
	}

	err = p.ws.WriteJSON(ANSWER, answer.SDP)
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
