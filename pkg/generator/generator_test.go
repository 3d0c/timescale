package generator

import (
	"encoding/csv"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestReadCSV(t *testing.T) {
	raw := `hostname,start_time,end_time
host_000008,2017-01-01 08:59:22,2017-01-01 09:59:22
host_000001,2017-01-02 13:02:02,2017-01-02 14:02:02
host_000008,2017-01-02 18:50:28,2017-01-02 19:50:28`

	expected := [][]string{
		{"hostname", "start_time", "end_time"},
		{"host_000008", "2017-01-01 08:59:22", "2017-01-01 09:59:22"},
		{"host_000001", "2017-01-02 13:02:02", "2017-01-02 14:02:02"},
		{"host_000008", "2017-01-02 18:50:28", "2017-01-02 19:50:28"},
	}

	csv := csv.NewReader(strings.NewReader(raw))
	if csv == nil {
		t.Fatalf("Error creating csv reader\n")
	}

	i := 0
	for {
		obtained, err := csv.Read()
		if err == io.EOF {
			break
		}

		if !reflect.DeepEqual(expected[i], obtained) {
			t.Fatalf("Expected %v, obtained %v\n", expected[i], obtained)
		}

		i++
	}
}
