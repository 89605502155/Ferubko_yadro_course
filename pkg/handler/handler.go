package handler

import (
	"net/http"

	"xkcd/pkg/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /", h.Auth)
	mux.HandleFunc("POST /update/", h.Update)
	mux.HandleFunc("GET /pics", h.Search)
	return mux

}
