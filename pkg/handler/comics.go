package handler

import (
	"fmt"
	"net/http"
)

func (h *Handler) Update(resp http.ResponseWriter, req *http.Request) {

	err := h.services.Comics.Update()
	fmt.Println(err)
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte("OK"))
	resp.Header().Set("Content-Type", "application/json")
}
