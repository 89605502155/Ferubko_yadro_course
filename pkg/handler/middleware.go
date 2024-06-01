package handler

import (
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

func Auth(f http.HandlerFunc, h *Handler, allowedSlice []string) http.HandlerFunc {
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
		userName, userStatus, err := h.services.Auth.ParseToken(headerParts[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
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
			http.Error(w, "You do not have the rights to perform this operation.",
				http.StatusForbidden)
			return
		}
		f(w, r)
		logrus.Println(userName, r.URL.Path)
	}
}
