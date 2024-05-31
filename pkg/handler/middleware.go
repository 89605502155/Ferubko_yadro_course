package handler

import "net/http"

func RateCheker(f http.HandlerFunc, h *Handler, hard int, dominantus bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rate_cheker := h.rate_limiter.Allow(hard, dominantus)
		if rate_cheker {
			f(w, r)
		} else {
			h.unavailableResponse(w)
			return
		}
	}
}
