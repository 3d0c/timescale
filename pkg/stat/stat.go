package stat

import (
	"fmt"
	"sync"
	"time"
)

type stat struct {
	sorted        []time.Duration
	totalNum      int
	totalDuration time.Duration
	avg           time.Duration
	distribution  map[int]int
	sync.Mutex
}

var (
	instance *stat
	once     sync.Once
)

func TheStat() *stat {
	once.Do(func() {
		instance = &stat{
			sorted:       make([]time.Duration, 0),
			distribution: make(map[int]int),
		}
	})

	return instance
}

func (s *stat) AddDuration(d time.Duration, wid int) {
	s.Lock()
	defer s.Unlock()

	s.totalNum++
	s.totalDuration += d
	s.distribution[wid]++

	if len(s.sorted) == 0 {
		s.sorted = append(s.sorted, d)
		return
	}

	if d > s.sorted[len(s.sorted)-1] {
		s.sorted = append(s.sorted, d)
	} else {
		s.sorted = append([]time.Duration{d}, s.sorted...)
	}

	return
}

func (s *stat) Dump() {
	s.Lock()
	defer s.Unlock()

	for _, d := range s.sorted {
		fmt.Printf("%v\n", d)
	}
}

func (s *stat) TotalNum() int {
	s.Lock()
	defer s.Unlock()

	result := s.totalNum
	return result
}

func (s *stat) TotalDuration() time.Duration {
	s.Lock()
	defer s.Unlock()

	result := s.totalDuration
	return result
}

func (s *stat) Distribution() {
	s.Lock()
	defer s.Unlock()

	keys := make([]int, 0)

	for key, _ := range s.distribution {
		keys = append(keys, key)
	}

	fmt.Printf("Distribution across workers:\n")

	for _, key := range keys {
		fmt.Printf("\tworker #%d%4d queries\n", key, s.distribution[key])
	}
}

func (s *stat) Min() time.Duration {
	s.Lock()
	defer s.Unlock()

	result := s.sorted[0]
	return result
}

func (s *stat) Max() time.Duration {
	s.Lock()
	defer s.Unlock()

	result := s.sorted[len(s.sorted)-1]
	return result
}

func (s *stat) Median() time.Duration {
	s.Lock()
	defer s.Unlock()

	result := s.sorted[len(s.sorted)/2]
	return result
}

func (s *stat) Avg() time.Duration {
	s.Lock()
	defer s.Unlock()

	result := time.Duration((s.sorted[0] + s.sorted[len(s.sorted)-1]) / 2)
	return result
}
