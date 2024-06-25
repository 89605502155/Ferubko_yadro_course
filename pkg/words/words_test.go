package words

import (
	"slices"
	"testing"
)

func compareMaps(map1, map2 *map[string]bool) bool {
	if len(*map1) != len(*map2) {
		return false
	}
	for key := range *map1 {
		if _, ok := (*map2)[key]; !ok {
			return false
		}
	}
	return true
}
func TestNormalization(t *testing.T) {
	testTable := []struct {
		sentence string
		expected *map[string]bool
	}{
		{
			sentence: "follower brings bunch of questions",
			expected: &map[string]bool{"bring": true, "bunch": true, "follow": true, "question": true},
		},
		{
			sentence: "i'll follow you as long as you are following me",
			expected: &map[string]bool{"follow": true, "ill": true, "long": true},
		},
		{
			sentence: "apple, doctor",
			expected: &map[string]bool{"appl": true, "doctor": true},
		},
		{
			sentence: "N-PLS1 is an algorithm of three-dimensional shape analysis",
			expected: &map[string]bool{"algorithm": true, "analysi": true, "npls": true, "shape": true,
				"threedimension": true},
		},
	}
	count := 0
	for _, testCase := range testTable {
		result, _ := NewStrimming().Normalization(testCase.sentence)
		if !compareMaps(result, testCase.expected) {
			t.Errorf("Expected %v, got %v", testCase.expected, result)
		} else {
			t.Logf("Expected  %v", testCase.expected)
			count++
		}
	}
	t.Logf("Work %d tests form %d", count, len(testTable))
}

func TestMergeMapToString(t *testing.T) {
	testTable := []struct {
		map1     *map[string]bool
		map2     *map[string]bool
		expected []string
	}{
		{
			map1:     &map[string]bool{"bring": true, "bunch": true, "follow": true, "question": true},
			map2:     &map[string]bool{"follow": true, "ill": true, "long": true},
			expected: []string{"bring", "bunch", "follow", "question", "ill", "long"},
		},
		{
			map1: &map[string]bool{"appl": true, "doctor": true},
			map2: &map[string]bool{"algorithm": true, "analysi": true, "npls": true, "shape": true,
				"threedimension": true},
			expected: []string{"appl", "doctor", "algorithm", "analysi", "npls", "shape", "threedimension"},
		},
	}
	count := 0
	for _, testCase := range testTable {
		result := NewStrimming().MergeMapToString(testCase.map1, testCase.map2)
		slices.Sort(result)
		slices.Sort(testCase.expected)
		if !slices.Equal(result, testCase.expected) {
			t.Errorf("Expected %v, got %v", testCase.expected, result)
		} else {
			t.Logf("Expected  %v", testCase.expected)
			count++
		}
	}
	t.Logf("Work %d tests form %d", count, len(testTable))
}
