package rate_limiter

import (
	"reflect"
	"testing"
	"time"
)

func TestNewSlidingLogLimiter(t *testing.T) {
	testTable := []struct {
		limit    int
		inteval  time.Duration
		expected *SlidindLogLimiter
	}{
		{
			10,
			time.Second,
			&SlidindLogLimiter{
				limit:   10,
				inteval: time.Second,
				logs:    make([]SlidingLog, 0),
			},
		},
	}
	count := 0
	for _, test := range testTable {
		res := NewSlidingLogLimiter(test.limit, test.inteval)
		if !reflect.DeepEqual(res, test.expected) {
			t.Errorf("Expected %v, got %v", test.expected, res)
		} else {
			t.Logf("Expected  %v", test.expected)
			count++
		}
	}
	t.Logf("Work %d tests form %d", count, len(testTable))
}

type rateQuery struct {
	hard       int
	moment     time.Duration
	dominantus bool
}

func TestAllow(t *testing.T) {
	long := make([]rateQuery, 1_000_000)
	// b1 := make([]bool, 1_000_000)
	b := make([]bool, 1_000_000)
	b2 := make([]bool, 1_000_000)
	// b1[0] = true

	for i := 0; i <= 9; i++ {
		b[i] = true
	}
	for i := range long {
		long[i].moment = 0
	}
	for i := range b2 {
		b2[i] = true
	}
	testTable := []struct {
		limit    int
		inteval  time.Duration
		data     []rateQuery
		expected []bool
	}{
		{
			10,
			time.Second,
			long,
			b,
		},
		{
			1e9,
			time.Microsecond,
			long,
			b2,
		},
	}
	for _, test := range testTable {
		sl := NewSlidingLogLimiter(test.limit, test.inteval)
		for i, ql := range test.data {
			if ql.moment != 0 {
				ticker := time.NewTicker(ql.moment * time.Microsecond)
				<-ticker.C
				ticker.Stop()
			}
			res := sl.Allow(ql.hard, ql.dominantus)
			if res != test.expected[i] {
				t.Errorf("Expected %v, got %v, %d", test.expected[i], res, i)
			}

		}
		t.Log("Next\n")
	}
}
