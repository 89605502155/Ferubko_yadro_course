package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
	"go.uber.org/mock/gomock"

	"xkcd/pkg/personal_limiter"
	"xkcd/pkg/rate_limiter"
	"xkcd/pkg/service"
	mock_service "xkcd/pkg/service/mocks"
)

func TestSearch(t *testing.T) {
	type mockBehavior func(s *mock_service.MockSearch, text string) string
	testTable := []struct {
		name               string
		inputBody          string
		mockBehavior       mockBehavior
		personalLimit      int
		personalInterval   time.Duration
		rateLimit          int
		rateInterval       time.Duration
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:      "first",
			inputBody: "follower brings bunch of questions",
			mockBehavior: func(s *mock_service.MockSearch, text string) string {
				var indexSlice, dbSlice []int
				var timeIndex, timeDb time.Duration
				var errorIndex, errorDb error
				s.EXPECT().SearchInDB(text).Return(dbSlice, timeDb, errorDb)
				s.EXPECT().SearchInIndex(text).Return(indexSlice, timeIndex, errorIndex)
				res := make(map[string]interface{})
				res["in db"] = dbSlice
				res["in index"] = indexSlice
				jn, _ := json.Marshal(res)
				return string(jn)
			},
			personalLimit:      5,
			personalInterval:   time.Second,
			rateLimit:          500,
			rateInterval:       time.Second,
			expectedStatusCode: 200,
			expectedBody:       "{\"in db\":[],\"in index\":[]}",
		},
	}
	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			search := mock_service.NewMockSearch(c)
			res1 := test.mockBehavior(search, test.inputBody)
			fmt.Println(res1)
			s := &service.Service{Search: search}
			p := personal_limiter.NewPersonalLimiter(context.Background(), test.personalLimit, test.personalInterval)
			r := rate_limiter.NewSlidingLogLimiter(test.rateLimit, test.rateInterval)
			mux := http.NewServeMux()
			mux.HandleFunc("GET /pics", NewHandler(s, r, p).Search)

			w := httptest.NewRecorder()
			queryParams := url.Values{}
			queryParams.Add("search", test.inputBody)
			req := httptest.NewRequest("GET", "/pics?"+queryParams.Encode(), nil) // Added query parameter

			mux.ServeHTTP(w, req)
			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, res1, w.Body.String())
		})
	}
}
