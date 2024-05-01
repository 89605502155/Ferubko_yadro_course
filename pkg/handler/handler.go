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
	// mux.HandleFunc("/" )
	mux.HandleFunc("/update/", h.Update)
	// mux.HandleFunc("/comic/random", h.GetRandomComic)
	// mux.HandleFunc("/comic/latest", h.GetLatestComic)
	// mux.HandleFunc("/comic/search", h.SearchComic)
	return mux

}
