package workerpool

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/3d0c/timescale/pkg/stat"
)

type Worker struct {
	id       int
	wChan    chan string
	hostname string
	db       *sql.DB
}

func NewWorker(hostname string, db *sql.DB, id int) *Worker {
	return &Worker{
		wChan:    make(chan string, 1),
		hostname: hostname,
		db:       db,
		id:       id,
	}
}

func (w *Worker) Id() int {
	return w.id
}

func (w *Worker) Send(query string) {
	w.wChan <- query
}

func (w *Worker) run(wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()

	for query := range w.wChan {
		start := time.Now()
		r, err := w.db.Query(query)
		stat.TheStat().AddDuration(time.Now().Sub(start), w.id)
		r.Close()

		if err != nil {
			fmt.Println(err)
		}
	}
}
