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
	searchAllowedSlice := []string{"user", "admin", "content manager"}
	updateAllowedSlice := []string{"admin", "content manager"}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /sign-in", RateCheker(h.SignIn, h, hardSearch, !dominantus))
	mux.HandleFunc("POST /create", h.CreateUser)
	mux.HandleFunc("POST /update/", RateCheker(Auth(h.Update, h, updateAllowedSlice), h, hard, dominantus))
	mux.HandleFunc("GET /pics", RateCheker(Auth(h.Search, h, searchAllowedSlice), h, hardSearch, !dominantus))
	return mux

}
