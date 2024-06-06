package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

func RateCheker(f http.HandlerFunc, h *Handler, hard int, dominantus bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rate_cheker := h.rate_limiter.Allow(hard, dominantus)
		if rate_cheker {
			f(w, r)
		} else {
			h.unavailableResponse(w)
			return
		}
	}
}

func (h *Handler) validToken(str string) (string, string, string, error) {
	userName, userStatus, err := h.services.Auth.ParseToken(str)
	logrus.Println("validToken ", err)
	if err != nil {
		logrus.Println("valid Token Refresh Token")
		userName, userStatus, access, err2 := h.services.Auth.ParseRefreshToken(str,
			accessTimeConst,
		)
		return userName, userStatus, access, err2
	}
	return userName, userStatus, str, err
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
		}
		notAllow := true
		for _, i := range allowedSlice {
			if i == userStatus {
				notAllow = false
				break
			}
		}
		if notAllow {
			http.Error(w, "You do not have the rights to perform this operation.",
				http.StatusForbidden)
			return
		}
		logic := h.personal_limiter.Allow(userName, hardTask)
		if logic {
			f(w, r)
		} else {
			http.Error(w, "Personal limiter.",
				http.StatusForbidden)
			return
		}
		logrus.Println(userName, r.URL.Path)
	}
}
