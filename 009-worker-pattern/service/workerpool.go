package service

type WorkerPool interface {
	Run()
	AddTask(task func())
}

type worker struct {
	poolSize        int
	queuedFunctions chan func()
}

func NewWorker(poolSize int) WorkerPool {
	return &worker{
		poolSize:        poolSize,
		queuedFunctions: make(chan func()),
	}
}

func (wr *worker) Run() {
	for i := 0; i < wr.poolSize; i++ {
		go func() {
			for task := range wr.queuedFunctions {
				task()
			}
		}()
	}
}

func (wr *worker) AddTask(task func()) {
	wr.queuedFunctions <- task
}
