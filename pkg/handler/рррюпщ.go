package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

type Handler struct {
	validToken       func(token string) (string, string, string, error)
	personal_limiter *PersonalLimiter
}

type PersonalLimiter struct {
	Allow func(userName string, hardTask int) bool
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func Auth(f http.HandlerFunc, h *Handler, allowedSlice []string, hardTask int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			http.Error(w, "Invalid auth header", http.StatusUnauthorized)
			return
		}
		userName, userStatus, tok, err := h.validToken(headerParts[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		logrus.Println(tok)
		logrus.Println(headerParts[1])
		if tok != headerParts[1] {
			tokens := TokenResponse{
				AccessToken:  tok,
				RefreshToken: headerParts[1],
			}
			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(tokens)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			logrus.Println("use refresh token")
			return
		}
		notAllow := true
		for _, i := range allowedSlice {
			if i == userStatus {
				notAllow = false
				break
			}
		}
		if notAllow {
			http.Error(w, "You do not have the rights to perform this operation.", http.StatusForbidden)
			return
		}
		logic := h.personal_limiter.Allow(userName, hardTask)
		if logic {
			f(w, r)
		} else {
			http.Error(w, "Personal limiter.", http.StatusForbidden)
			return
		}
		logrus.Println(userName, r.URL.Path)
	}
}

func TestAuth(t *testing.T) {
	handler := &Handler{
		validToken: func(token string) (string, string, string, error) {
			if token == "valid-token" {
				return "user1", "admin", "valid-token", nil
			}
			return "", "", "", http.ErrNoLocation
		},
		personal_limiter: &PersonalLimiter{
			Allow: func(userName string, hardTask int) bool {
				return true
			},
		},
	}

	allowedSlice := []string{"admin", "user"}

	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}

	wrappedHandler := Auth(testHandler, handler, allowedSlice, 10)

	tests := []struct {
		name       string
		authHeader string
		wantStatus int
		wantBody   string
	}{
		{"NoAuthHeader", "", http.StatusUnauthorized, "Authorization header is missing\n"},
		{"InvalidAuthHeader", "InvalidHeader", http.StatusUnauthorized, "Invalid auth header\n"},
		{"InvalidToken", "Bearer invalid-token", http.StatusUnauthorized, "no Location header in response\n"},
		{"ValidToken", "Bearer valid-token", http.StatusOK, "OK"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rr := httptest.NewRecorder()
			wrappedHandler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
			if body := rr.Body.String(); body != tt.wantBody {
				t.Errorf("handler returned unexpected body: got %v want %v", body, tt.wantBody)
			}
		})
	}
}
