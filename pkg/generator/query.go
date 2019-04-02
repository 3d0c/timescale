package generator

import (
	"fmt"
	"time"
)

type Record struct {
	hostname string
	start    time.Time
	end      time.Time
}

func NewRecord(host, startTime, endTime string) (*Record, error) {
	var err error

	r := &Record{hostname: host}

	if r.start, err = time.Parse("2006-01-02 15:04:05", startTime); err != nil {
		return nil, fmt.Errorf("error parsing '%s' - %s\n", startTime, err)
	}

	if r.end, err = time.Parse("2006-01-02 15:04:05", endTime); err != nil {
		return nil, fmt.Errorf("error parsing '%s' - %s\n", endTime, err)
	}

	return r, nil
}

func (r *Record) minutes() []time.Time {
	result := make([]time.Time, 0)

	start := r.start

	for {
		if start.After(r.end) {
			break
		}

		result = append(result, start)
		start = start.Add(time.Minute)
	}

	return result
}

func (r *Record) Queries() chan Query {
	yield := make(chan Query)
	go func() {
		defer close(yield)
		for _, ts := range r.minutes() {
			q := Query{
				hostname:  r.hostname,
				timestamp: ts,
			}

			yield <- q
		}
	}()

	return yield
}

type Query struct {
	hostname  string
	timestamp time.Time
}

func (q *Query) String() string {
	return fmt.Sprintf("select min(usage), max(usage) from cpu_usage where to_char(ts, 'YYYY-MM-DD HH24:MI') = '%s' and host = '%s'", q.timestamp.Format("2006-01-02 15:04"), q.hostname)
}

func (q *Query) GetHostName() string {
	return q.hostname
}
