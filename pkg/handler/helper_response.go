package handler

import "net/http"

func (h *Handler) unavailableResponse(resp http.ResponseWriter) {
	resp.WriteHeader(http.StatusServiceUnavailable)
	resp.Header().Set("Content-Type", "application/json")
	resp.Header().Set("Retry-After", "600")
	resp.Write([]byte("Сервер занят другим запросом"))
}

func (h *Handler) okStatusResponse(resp http.ResponseWriter) {
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte("OK"))
	resp.Header().Set("Content-Type", "application/json")
}
