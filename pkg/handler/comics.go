package handler

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

const (
	hard       = 250_000
	dominantus = true
)

func (h *Handler) Update(resp http.ResponseWriter, req *http.Request) {
	err := h.services.Comics.Update()
	if err != nil {
		logrus.Error(err)
		resp.WriteHeader(http.StatusInternalServerError)
	} else {
		h.okStatusResponse(resp)
	}
}
