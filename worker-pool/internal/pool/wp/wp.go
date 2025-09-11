package wp

import "fmt"

type Worker struct {
	id               int
	handledJobsCount int
}

var workers = []*Worker{}

func InitWorkers(workersCnt int) {
	for i := range workersCnt {
		workers = append(workers, &Worker{id: i + 1})
	}
}

type Pool[Data any] struct {
	pool    chan *Worker
	handler func(int, Data)
}

func New[Data any](handler func(int, Data)) *Pool[Data] {
	return &Pool[Data]{
		handler: handler,
		pool:    make(chan *Worker, len(workers)),
	}
}

func (p *Pool[Data]) Create() {
	for _, w := range workers {
		p.pool <- w
	}
}
func (p *Pool[Data]) Handle(data Data) {
	worker := <-p.pool

	go func() {
		p.handler(worker.id, data)
		worker.handledJobsCount++
		p.pool <- worker
	}()

}
func (p *Pool[Data]) Wait() {

	for range len(workers) {
		<-p.pool
	}
}

func (p *Pool[Data]) Stats() {
	fmt.Println("__________STATS________")
	for _, worker := range workers {
		fmt.Printf("worker: %d, total Jobs: %d\n", worker.id, worker.handledJobsCount)
	}
}
