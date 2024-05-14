package service

import (
	"sync"
)

type WorkerPool interface {
	Run()
	AddTask(task func() error)
	Wait()
}

type worker struct {
	poolSize    int
	queuedTasks chan func() error
	wg          *sync.WaitGroup
	stopOnce    *sync.Once
}

func NewWorker(poolSize int) WorkerPool {
	return &worker{
		poolSize:    poolSize,
		queuedTasks: make(chan func() error),
		wg:          &sync.WaitGroup{},
		stopOnce:    &sync.Once{},
	}
}

func (wr *worker) Run() {
	for i := 0; i < wr.poolSize; i++ {
		go func() {
			var err error
			for task := range wr.queuedTasks {
				if err == nil {
					err = task()
				}
				wr.wg.Done()
			}
		}()
	}
}

func (wr *worker) AddTask(task func() error) {
	wr.wg.Add(1)
	wr.queuedTasks <- task
}

func (wr *worker) Wait() {
	wr.stopOnce.Do(func() {
		wr.wg.Wait()
		close(wr.queuedTasks)
	})
}
