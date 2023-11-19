package handler

import "sync"

type WorkerHandler struct {
	wg sync.WaitGroup
	ch chan struct{}
}

func GetWorkerHandler(workerCount int) *WorkerHandler {
	return &WorkerHandler{
		wg: sync.WaitGroup{},
		ch: make(chan struct{}, workerCount),
	}
}

func (h *WorkerHandler) Run(i int, f func(int)) {
	h.ch <- struct{}{}
	h.wg.Add(1)
	go func() {
		defer func() {
			h.wg.Done()
			<-h.ch
		}()
		f(i)
	}()
}

func (h *WorkerHandler) Wait() {
	h.wg.Wait()
}
