package workerpool

import (
	"database/sql"
	"math/rand"
	"sync"
	"time"
)

type Pool struct {
	num     int
	wg      *sync.WaitGroup
	workers []*Worker
}

func NewPool(wnum int, db *sql.DB) *Pool {
	pool := &Pool{
		num:     wnum,
		wg:      &sync.WaitGroup{},
		workers: make([]*Worker, 0),
	}

	for i := 0; i < pool.num; i++ {
		w := NewWorker("", db, i)
		pool.wg.Add(1)
		go w.run(pool.wg)
		pool.workers = append(pool.workers, w)
	}

	return pool
}

func (p *Pool) Get(name string) (w *Worker) {
	if w = p.getWorker(name); w == nil {
		w = p.getRandomWorker()
	}

	if w == nil {
		panic("unexpected error")
	}

	return w
}

func (p *Pool) StopWorkers() {
	for i := 0; i < p.num; i++ {
		close(p.workers[i].wChan)
	}
}

func (p *Pool) Wait() {
	p.wg.Wait()
}

func (p *Pool) getWorker(hostname string) (w *Worker) {
	for _, w = range p.workers {
		if w.hostname == hostname {
			return w
		}
		if w.hostname == "" {
			w.hostname = hostname
			return w
		}
	}

	return nil
}

func (p *Pool) getRandomWorker() *Worker {
	rand.Seed(time.Now().UnixNano())
	return p.workers[rand.Intn(len(p.workers))]
}
