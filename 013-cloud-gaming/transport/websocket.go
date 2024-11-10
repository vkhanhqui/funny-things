package transport

import (
	"cloud-gaming/service"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v4"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
)

const STUN = "stun:stun.l.google.com:19302"

func NewHandler(svc service.Service) *Handler {
	return &Handler{svc}
}

type Handler struct {
	svc service.Service
}

func (h *Handler) Websocket(w http.ResponseWriter, r *http.Request) {
	unsafeConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	conn := service.NewSafeWriter(unsafeConn)
	pc, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{STUN},
			},
		},
	})
	if err != nil {
		conn.WriteJSON(service.ERROR, "Error creating peer connection")
		return
	}

	p := service.NewPeerConnection(pc, conn)
	go h.svc.HandleMessage(p)
}
