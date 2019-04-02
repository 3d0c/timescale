package workerpool

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	var (
		cases []string = []string{"1", "2", "3", "4"}
		w     *worker
	)

	pool := NewPool(2)

	for i, hostname := range cases {
		if w = pool.getWorker(hostname); w == nil {
			if w = pool.getRandomWorker(); w == nil {
				t.Fatalf("Unexpected error\n")
			}
		}

		if (w.hostname != hostname) && (i != 2 && i != 3) {
			t.Fatalf("Expected hostname = %s, obtained = %s\n", hostname, w.hostname)
		}
	}
}

func TestRunPool(t *testing.T) {
	pool := NewPool(10)

	go func() {
		time.Sleep(time.Second * 3)
		pool.StopWorkers()
	}()

	for i := 0; i < 20; i++ {
		w := pool.Get(strconv.Itoa(i))
		w.Send(fmt.Sprintf("hostname %d", i))
	}

	pool.Wait()
}
