package generator

import (
	"fmt"
	"strconv"
	"testing"
)

func TestTimeRange(t *testing.T) {
	r, err := NewRecord("host", "2017-01-01 08:59:22", "2017-01-01 09:59:22")
	if err != nil {
		t.Fatalf("%s\n", err)
	}

	result := r.minutes()
	if n := len(result); n != 61 {
		t.Fatalf("Expected result length = 62, obtained - %d", n)
	}
}

func TestQuery(t *testing.T) {
	r, err := NewRecord("host", "2017-01-01 08:50:22", "2017-01-01 08:51:22")
	if err != nil {
		t.Fatalf("%s\n", err)
	}

	result := make([]Query, 0)

	for i, ts := range r.minutes() {
		result = append(result, Query{hostname: "host" + strconv.Itoa(i), timestamp: ts})
	}

	if n := len(result); n != 2 {
		for _, r := range result {
			fmt.Println(r.hostname, r.String())
		}

		t.Fatalf("Expected result length = 2, obtained - %d", n)
	}

	expected := []string{
		"select min(usage), max(usage) from cpu_usage where to_char(ts, 'YYYY-MM-DD HH24:MI') = '2017-01-01 08:50' and host = 'host0'",
		"select min(usage), max(usage) from cpu_usage where to_char(ts, 'YYYY-MM-DD HH24:MI') = '2017-01-01 08:51' and host = 'host1'",
	}

	for i, r := range result {
		if expected[i] != r.String() {
			t.Fatalf("Expected '%s', obtained '%s'\n", expected[i], r.String())
		}
	}
}
