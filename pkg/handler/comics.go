package handler

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func (h *Handler) Auth(resp http.ResponseWriter, req *http.Request) {
	logrus.Println("Auth", req.Method)
	resp.Write([]byte("You are user\n"))
	resp.Header().Set("Token", "<PASSWORD>")
}

func (h *Handler) Update(resp http.ResponseWriter, req *http.Request) {
	h.Auth(resp, req)
	err := h.services.Comics.Update()
	logrus.Println(err)
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte("OK"))
	resp.Header().Set("Content-Type", "application/json")
}
