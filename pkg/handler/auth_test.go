package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
	"go.uber.org/mock/gomock"

	server "xkcd"
	"xkcd/pkg/personal_limiter"
	"xkcd/pkg/rate_limiter"
	"xkcd/pkg/service"
	mock_service "xkcd/pkg/service/mocks"
)

func TestSignIn(t *testing.T) {

	type mockBehavior func(s *mock_service.MockAuth, data server.UserEntity)
	testTable := []struct {
		name               string
		inputBody          string
		data               server.UserEntity
		mockBehavior       mockBehavior
		personalLimit      int
		personalInterval   time.Duration
		rateLimit          int
		rateInterval       time.Duration
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:      "success",
			inputBody: `{"username": "andrey", "password": "123"}`,
			data: server.UserEntity{
				Username: "andrey",
				Password: "123",
			},
			mockBehavior: func(s *mock_service.MockAuth, data server.UserEntity) {
				if err := data.Validate(); err == nil {
					s.EXPECT().GenerateToken(data, accessTimeConst, refreshTimeConst).Return("accessToken", "refreshToken", err)
				}
			},
			personalLimit:      5,
			personalInterval:   time.Second,
			rateLimit:          500,
			rateInterval:       time.Second,
			expectedStatusCode: 200,
			expectedBody:       `{"accessToken":"accessToken","refreshToken":"refreshToken"}`,
		},
		{
			name:      "no-password",
			inputBody: `{"username": "andrey", "password": ""}`,
			data: server.UserEntity{
				Username: "andrey",
				Password: "",
			},
			mockBehavior: func(s *mock_service.MockAuth, data server.UserEntity) {
				if err := data.Validate(); err == nil {
					s.EXPECT().GenerateToken(data, accessTimeConst, refreshTimeConst).Return("accessToken", "refreshToken", err)
				}
			},
			personalLimit:      5,
			personalInterval:   time.Second,
			rateLimit:          500,
			rateInterval:       time.Second,
			expectedStatusCode: 400,
			expectedBody:       "",
		},
		{
			name:      "no-name",
			inputBody: `{"username": "", "password": "123"}`,
			data: server.UserEntity{
				Username: "",
				Password: "123",
			},
			mockBehavior: func(s *mock_service.MockAuth, data server.UserEntity) {
				if err := data.Validate(); err == nil {
					s.EXPECT().GenerateToken(data, accessTimeConst, refreshTimeConst).Return("accessToken", "refreshToken", err)
				}
			},
			personalLimit:      5,
			personalInterval:   time.Second,
			rateLimit:          500,
			rateInterval:       time.Second,
			expectedStatusCode: 400,
			expectedBody:       "",
		},
		{
			name:      "invalid password",
			inputBody: `{"username": "andrey", "password": "wrongpassword"}`,
			data: server.UserEntity{
				Username: "andrey",
				Password: "wrongpassword",
			},
			mockBehavior: func(s *mock_service.MockAuth, data server.UserEntity) {
				if err := data.Validate(); err == nil {
					s.EXPECT().GenerateToken(data, accessTimeConst, refreshTimeConst).Return("", "", errors.New("invalid password"))
				}
			},
			personalLimit:      5,
			personalInterval:   time.Second,
			rateLimit:          500,
			rateInterval:       time.Second,
			expectedStatusCode: 401,
			expectedBody:       "",
		},
	}
	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuth(c)
			test.mockBehavior(auth, test.data)

			s := &service.Service{Auth: auth}
			p := personal_limiter.NewPersonalLimiter(context.Background(), test.personalLimit, test.personalInterval)
			r := rate_limiter.NewSlidingLogLimiter(test.rateLimit, test.rateInterval)
			mux := http.NewServeMux()
			mux.HandleFunc("POST /sign-in", RateCheker(NewHandler(s, r, p).SignIn, NewHandler(s, r, p), hardSearch, !dominantus))

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-in", bytes.NewBufferString(test.inputBody))

			mux.ServeHTTP(w, req)
			assert.Equal(t, test.expectedStatusCode, w.Code)

			var expectedBody, actualBody map[string]interface{}
			_ = json.Unmarshal([]byte(test.expectedBody), &expectedBody)
			_ = json.Unmarshal(w.Body.Bytes(), &actualBody)

			assert.Equal(t, expectedBody, actualBody)

		})
	}
}
