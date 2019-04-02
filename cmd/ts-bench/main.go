package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/lib/pq"

	"github.com/3d0c/timescale/pkg/generator"
	"github.com/3d0c/timescale/pkg/stat"
	"github.com/3d0c/timescale/pkg/workerpool"
)

func main() {
	var (
		fp         *os.File
		db         *sql.DB
		dbargs     string
		paramsFile string
		wnum       int
		err        error
		gen        *generator.Generator
	)

	flag.StringVar(&paramsFile, "qp", "", "Query parameters file")
	flag.IntVar(&wnum, "wnum", 4, "Number of workers")
	flag.StringVar(&dbargs, "dbargs", "postgres://alex:@localhost/homework?sslmode=disable", "Postgres connection arguments")
	flag.Parse()

	if db, err = sql.Open("postgres", dbargs); err != nil {
		log.Fatalf("Error opening postgres - %s\n", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Error connection to postgres - %s\n", err)
	}

	stat.TheStat()

	if fp, err = os.OpenFile(paramsFile, os.O_RDONLY, 0444); err != nil {
		log.Fatalf("Error opening parameters file - %s\n", err)
	}
	defer fp.Close()

	pool := workerpool.NewPool(wnum, db)

	if gen, err = generator.New(fp); err != nil {
		log.Fatalf("Error initializing query generator - %s\n", err)
	}

	muxq := make(chan generator.Query)
	muxw := make(chan *workerpool.Worker)

	wg := sync.WaitGroup{}

	for i := 0; i < wnum; i++ {
		wg.Add(1)
		go func(chan *workerpool.Worker, chan generator.Query) {
			defer wg.Done()
			for {
				w, ok := <-muxw
				if !ok {
					return
				}
				q, ok := <-muxq
				if !ok {
					return
				}
				w.Send(q.String())
			}
		}(muxw, muxq)
	}

	for query := range gen.Query() {
		muxw <- pool.Get(query.GetHostName())
		muxq <- query
	}

	close(muxw)
	close(muxq)

	wg.Wait()

	pool.StopWorkers()
	pool.Wait()

	fmt.Printf("Total queries: %d\n", stat.TheStat().TotalNum())
	stat.TheStat().Distribution()
	fmtstr := "\t%-8s %-4v\n"
	fmt.Printf("Duration:\n")
	fmt.Printf(fmtstr, "total:", stat.TheStat().TotalDuration())
	fmt.Printf(fmtstr, "minimum:", stat.TheStat().Min())
	fmt.Printf(fmtstr, "maximum:", stat.TheStat().Max())
	fmt.Printf(fmtstr, "median:", stat.TheStat().Median())
	fmt.Printf(fmtstr, "average:", stat.TheStat().Avg())
}
