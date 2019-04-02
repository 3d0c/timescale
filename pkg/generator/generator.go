package generator

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
)

type Generator struct {
	csv *csv.Reader
}

func New(r io.Reader) (*Generator, error) {
	g := &Generator{}

	if g.csv = csv.NewReader(r); g.csv == nil {
		return nil, fmt.Errorf("error creating csv.NewReader\n")
	}

	return g, nil
}

func (g *Generator) Query() chan Query {
	yield := make(chan Query)

	go func() {
		var (
			line  int
			rcnum int
		)

		defer close(yield)

		for {
			items, err := g.csv.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Printf("Error reading line %d - %s\n", line, err)
				continue
			}

			line++

			if line == 1 {
				rcnum = len(items)
				continue
			}

			if len(items) != rcnum {
				log.Printf("Malformed CSV entry - %v\n", items)
				continue
			}

			r, err := NewRecord(items[0], items[1], items[2])
			if err != nil {
				log.Printf("Error creating new record - malformed CSV entry - %v\n", items)
				continue
			}

			for query := range r.Queries() {
				yield <- query
			}
		}
	}()

	return yield
}
