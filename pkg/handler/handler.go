package handler

import (
	"net/http"

	"xkcd/pkg/rate_limiter"
	"xkcd/pkg/service"
)

type Handler struct {
	services     *service.Service
	rate_limiter *rate_limiter.SlidindLogLimiter
}

func NewHandler(services *service.Service, rate_limiter *rate_limiter.SlidindLogLimiter) *Handler {
	return &Handler{
		services:     services,
		rate_limiter: rate_limiter,
	}
}

func (h *Handler) InitRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /", h.Auth)
	mux.HandleFunc("POST /update/", h.Update)
	mux.HandleFunc("GET /pics", h.Search)
	return mux

}
