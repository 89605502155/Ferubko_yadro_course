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
				limit:          10,
				inteval:        time.Second,
				logs:           make([]SlidingLog, 0),
				timeDominantus: time.Now().Add(-timeBetween),
			},
		},
	}
	count := 0
	for _, test := range testTable {
		res := NewSlidingLogLimiter(test.limit, test.inteval)
		rep := struct {
			limit   int
			inteval time.Duration
			logs    []SlidingLog
		}{
			res.limit, res.inteval, res.logs,
		}
		tres := struct {
			limit   int
			inteval time.Duration
			logs    []SlidingLog
		}{
			test.expected.limit, test.expected.inteval, test.expected.logs,
		}
		if !reflect.DeepEqual(rep, tres) {
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
	b := make([]bool, 1_000_000)

	for i := 0; i <= 9; i++ {
		b[i] = true
	}
	for i := range long {
		long[i].moment = 0
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
			5,
			time.Second,
			[]rateQuery{
				{0, 0, false}, {500, 1, true}, {0, 0, false},
			},
			[]bool{true, true, false},
		},
		{
			5,
			time.Second,
			[]rateQuery{
				{0, 0, false}, {500_000, 1, true}, {0, 0, false},
			},
			[]bool{true, true, false},
		},
		{
			5,
			time.Second,
			[]rateQuery{
				{0, 0, false}, {500, 1, true}, {500, 0, true},
			},
			[]bool{true, true, false},
		},
		{
			5,
			time.Second,
			[]rateQuery{
				{0, 0, false}, {500, 1, true}, {500, 5_000_000, true},
			},
			[]bool{true, true, false},
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
