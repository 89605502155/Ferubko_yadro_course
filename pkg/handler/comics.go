package handler

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

const (
	hard       = 250
	dominantus = true
)

func (h *Handler) Auth(resp http.ResponseWriter, req *http.Request) {
	logrus.Println("Auth", req.Method)
	resp.Write([]byte("You are user\n"))
	resp.Header().Set("Token", "<PASSWORD>")
}

func (h *Handler) Update(resp http.ResponseWriter, req *http.Request) {
	rate_cheker := h.rate_limiter.Allow(hard, dominantus)
	if rate_cheker {
		h.Auth(resp, req)
		err := h.services.Comics.Update()
		if err != nil {
			logrus.Error(err)
			resp.WriteHeader(http.StatusInternalServerError)
		} else {
			h.okStatusResponse(resp)
		}
	} else {
		h.unavailableResponse(resp)
	}
}
