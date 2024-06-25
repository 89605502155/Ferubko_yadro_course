package handler

import (
	"net/http"

	"xkcd/pkg/personal_limiter"
	"xkcd/pkg/rate_limiter"
	"xkcd/pkg/service"
)

type Handler struct {
	services         *service.Service
	rate_limiter     *rate_limiter.SlidindLogLimiter
	personal_limiter *personal_limiter.PersonalLimiter
}

func NewHandler(services *service.Service, rate_limiter *rate_limiter.SlidindLogLimiter,
	personal_limiter *personal_limiter.PersonalLimiter) *Handler {
	return &Handler{
		services:         services,
		rate_limiter:     rate_limiter,
		personal_limiter: personal_limiter,
	}
}

func (h *Handler) InitRoutes() http.Handler {
	searchAllowedSlice := []string{"user", "admin", "content manager"}
	updateAllowedSlice := []string{"admin", "content manager"}
	createUserSlice := []string{"admin"}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /sign-in", RateCheker(h.SignIn, h, hardSearch, !dominantus))
	mux.HandleFunc("POST /create", RateCheker(Auth(h.CreateUser, h, createUserSlice, hardSearch), h, hardSearch, !dominantus))
	mux.HandleFunc("POST /update/", RateCheker(Auth(h.Update, h, updateAllowedSlice, hard), h, hard, dominantus))
	mux.HandleFunc("GET /pics", RateCheker(Auth(h.Search, h, searchAllowedSlice, hardSearch), h, hardSearch, !dominantus))
	return mux

}
