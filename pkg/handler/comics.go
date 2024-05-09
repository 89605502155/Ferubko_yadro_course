package handler

import (
	"fmt"
	"net/http"
)

func (h *Handler) Auth(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("Auth", req.Method)
	// resp.WriteHeader(http.StatusOK)
	resp.Write([]byte("You are user\n"))
	resp.Header().Set("Token", "<PASSWORD>")
	// resp.Header().Set("Content-Type", "application/json")
}

func (h *Handler) Update(resp http.ResponseWriter, req *http.Request) {
	h.Auth(resp, req)
	err := h.services.Comics.Update()
	fmt.Println(err)
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte("OK"))
	resp.Header().Set("Content-Type", "application/json")
}
