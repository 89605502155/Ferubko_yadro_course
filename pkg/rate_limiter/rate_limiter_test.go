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
