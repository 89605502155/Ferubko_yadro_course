package personal_limiter

import (
	"context"
	"reflect"
	"testing"
	"time"
)


func TestPersonalLimiter(t *testing.T) {
	testTable := []struct {
		ctx      context.Context
		limit    int
		interval time.Duration
		expected *PersonalLimiter
	}{
		{
			context.Background(),
			10,
			time.Second,
			&PersonalLimiter{
				limit:   10,
				inteval: time.Second,
				list:    make(map[string][]int, 0),
			},
		},
	}
	count := 0
	for _, testCase := range testTable {
		result := NewPersonalLimiter(testCase.ctx, testCase.limit, testCase.interval)
		if !reflect.DeepEqual(result, testCase.expected) {
			t.Errorf("Expected %v, got %v", testCase.expected, result)
		} else {
			t.Logf("Expected  %v", testCase.expected)
			count++
		}
	}
	t.Logf("Work %d tests form %d", count, len(testTable))
}

type queries struct {
	userName string
	hard     int
	moment   time.Duration
}

func TestPersonAllow(t *testing.T) {
	testTable := []struct {
		ctx      context.Context
		limit    int
		interval time.Duration
		data     []queries
		expected []bool
	}{
		{
			context.Background(),
			10,
			time.Second,
			[]queries{
				{"user", 9, 0}, {"user", 10, 0}, {"user", 10, 0},
				{"user", 10, 0}, {"user", 10, 0}, {"user", 10, 0},
				{"user", 10, 0}, {"admin", 500, 0}, {"admin", 5, 0},
				{"a", 5, 0}, {"admin", 5, 0},
			},
			[]bool{true, false, false, false, false, false, false, false, true, true, false},
		},
		{
			context.Background(),
			900,
			time.Second,
			[]queries{
				{"user", 9, 0}, {"user", 10, 0}, {"user", 10, 0},
				{"user", 10, 0}, {"user", 10, 0}, {"user", 10, 0},
				{"user", 10, 0}, {"admin", 500, 0}, {"admin", 5, 0},
				{"a", 5, 0}, {"admin", 5, 0},
			},
			[]bool{true, true, true, true, true, true, true, true, true, true, true},
		},
		{
			context.Background(),
			10,
			time.Second,
			[]queries{
				{"user", 9, 0}, {"user", 10, 0}, {"user", 10, 0},
				{"user", 10, 0}, {"user", 10, 0}, {"user", 10, 0},
				{"user", 10, 0}, {"admin", 500, 0}, {"admin", 5, 0},
				{"a", 5, 0}, {"admin", 3, 0},
			},
			[]bool{true, false, false, false, false, false, false, false, true, true, true},
		},
		{
			context.Background(),
			1,
			time.Second,
			[]queries{
				{"user", 0, 0}, {"user", 0, 0}, {"user", 0, 0},
				{"user", 10, 0}, {"user", 10, 0}, {"user", 10, 0},
				{"user", 10, 0}, {"admin", 500, 0}, {"admin", 0, 0},
				{"a", 0, 0}, {"admin", 0, 0},
			},
			[]bool{true, false, false, false, false, false, false, false, true, true, false},
		},
		{
			context.Background(),
			10,
			time.Second,
			[]queries{
				{"user", 9, 0}, {"user", 10, 1000}, {"user", 10, 1000},
				{"user", 10, 0}, {"user", 10, 0}, {"user", 10, 0},
				{"user", 10, 0}, {"admin", 500, 10}, {"admin", 5, 1200},
				{"a", 5, 0}, {"admin", 5, 0},
			},
			[]bool{true, false, false, false, false, false, false, false, true, true, false},
		},
		{
			context.Background(),
			10,
			time.Second,
			[]queries{
				{"user", 9, 0}, {"user", 10, 1100}, {"user", 10, 100},
				{"user", 10, 0}, {"user", 10, 0}, {"user", 10, 0},
				{"user", 10, 0}, {"admin", 500, 10}, {"admin", 5, 0},
				{"a", 5, 1200}, {"admin", 3, 3900}, {"a", 6, 1900},
			},
			[]bool{true, false, false, false, false, false, false, false, true, true, true, true},
		},
	}

	for _, testCase := range testTable {
		limitStruct := NewPersonalLimiter(testCase.ctx, testCase.limit, testCase.interval)
		for i, data := range testCase.data {
			if data.moment != 0 {
				time.Sleep(data.moment * time.Millisecond)
			}
			res := limitStruct.Allow(data.userName, data.hard)
			if res != testCase.expected[i] {
				t.Errorf("Expected %v, got %v in %d", testCase.expected[i], res, i)
			}
		}
		t.Logf("Querry map %v", limitStruct.list)
	}
}
