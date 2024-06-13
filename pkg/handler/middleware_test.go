package handler

import (
	"fmt"
	"net/http"
	"testing"
)

func TestRateCheker(t *testing.T) {
	testTable := []struct {
		hard       int
		dominant   bool
		statusCode int
	}{
		{0, true, 200},
	}
	for _, test := range testTable {
		var f http.HandlerFunc
		// var w http.ResponseWriter
		// var r *http.Request
		h := &Handler{}
		res := RateCheker(f, h, test.hard, test.dominant)
		fmt.Println(res)

	}
}

func TestAuth(t *testing.T) {

}
