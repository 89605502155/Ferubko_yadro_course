package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
	"go.uber.org/mock/gomock"

	"xkcd/pkg/personal_limiter"
	"xkcd/pkg/rate_limiter"
	"xkcd/pkg/service"
	mock_service "xkcd/pkg/service/mocks"
)

func TestUpdate(t *testing.T) {
	type mockBehavior func(s *mock_service.MockComics)
	testTable := []struct {
		name               string
		mockBehavior       mockBehavior
		personalLimit      int
		personalInterval   time.Duration
		rateLimit          int
		rateInterval       time.Duration
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name: "success",
			mockBehavior: func(s *mock_service.MockComics) {
				s.EXPECT().Update().Return(nil)
			},
			personalLimit:      5,
			personalInterval:   time.Second,
			rateLimit:          500,
			rateInterval:       time.Second,
			expectedStatusCode: 200,
			expectedBody:       "OK",
		},
	}
	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			comic := mock_service.NewMockComics(c)
			test.mockBehavior(comic)

			s := &service.Service{Comics: comic}
			p := personal_limiter.NewPersonalLimiter(context.Background(), test.personalLimit, test.personalInterval)
			r := rate_limiter.NewSlidingLogLimiter(test.rateLimit, test.rateInterval)
			mux := http.NewServeMux()
			mux.HandleFunc("POST /update/", NewHandler(s, r, p).Update)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/update/", nil)

			mux.ServeHTTP(w, req)
			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedBody, w.Body.String())
		})
	}
}
