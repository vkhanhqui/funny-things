package service

import "log"

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
		}
	}()
}
