package transport

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v4"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	// peerConns  []*PeerConnState
	peerConnCh = make(chan *PeerConnState)
)

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	unsafeConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	conn := &threadSafeWriter{unsafeConn, sync.Mutex{}}
	pc, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	})
	if err != nil {
		conn.WriteJSON(Message{Type: "error", Value: "Error creating peer connection"})
		return
	}

	pcs := &PeerConnState{peerConnection: pc, ws: conn}
	pc.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		log.Println("Connection state changed:", state)
		if state == webrtc.PeerConnectionStateConnected {
			pcs = &PeerConnState{peerConnection: pc, ws: conn}
			peerConnCh <- pcs
			log.Println("Added new peer connection to channel")
		}
	})

	go handleMessage(pcs)
}

func handleMessage(pcs *PeerConnState) {
	for {
		var msg Message
		err := pcs.ws.ReadJSON(&msg)
		if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseAbnormalClosure) {
			log.Printf("Websocket closed: %v", err)
			break
		} else if err != nil {
			log.Printf("Error reading from websocket: %v", err)
			continue
		}

		switch msg.Type {
		case "offer":
			offer := msg.Value.(string)
			err = pcs.peerConnection.SetRemoteDescription(webrtc.SessionDescription{Type: webrtc.SDPTypeOffer, SDP: offer})
			if err != nil {
				log.Printf("Error reading from websocket: %v", err)
				continue
			}

			answer, err := pcs.peerConnection.CreateAnswer(nil)
			if err != nil {
				log.Printf("Error reading from websocket: %v", err)
				continue
			}

			err = pcs.peerConnection.SetLocalDescription(answer)
			if err != nil {
				log.Printf("Error reading from websocket: %v", err)
				continue
			}
			pcs.ws.WriteJSON(Message{Type: "answer", Value: answer.SDP})

		case "ice":
			ice := msg.Value.(string)
			pcs.peerConnection.AddICECandidate(webrtc.ICECandidateInit{Candidate: ice})

		default:
			log.Printf("Unknown message type: %+v", msg)
			continue
		}
	}
}

type Message struct {
	Type  string `json:"type"`
	Value any    `json:"value"`
}

type PeerConnState struct {
	peerConnection *webrtc.PeerConnection
	ws             *threadSafeWriter
}

type threadSafeWriter struct {
	*websocket.Conn
	sync.Mutex
}

func (t *threadSafeWriter) WriteJSON(v interface{}) error {
	t.Lock()
	defer t.Unlock()
	return t.Conn.WriteJSON(v)
}
