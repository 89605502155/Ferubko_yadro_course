package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type Handler struct {
	rate_limiter        *RateLimiter
	unavailableResponse func(w http.ResponseWriter)
}

type RateLimiter struct {
	// Здесь можно добавить любые поля, необходимые для ограничения скорости
}

func (rl *RateLimiter) Allow(hard int, dominantus bool) bool {
	// Здесь будет логика для ограничения скорости
	return true
}

func RateChecker(f http.HandlerFunc, h *Handler, hard int, dominantus bool) http.HandlerFunc {
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

func TestRateChecker(t *testing.T) {
	// Создаем фейковый RateLimiter, который всегда возвращает true
	rateLimiter := &RateLimiter{}
	handler := &Handler{
		rate_limiter: rateLimiter,
		unavailableResponse: func(w http.ResponseWriter) {
			http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		},
	}

	// Создаем функцию, которая будет обернута в RateChecker
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}

	// Оборачиваем тестовую функцию в RateChecker
	wrappedHandler := RateChecker(testHandler, handler, 10, true)

	// Создаем новый запрос и запись для ответа
	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	rr := httptest.NewRecorder()

	// Вызываем обернутый обработчик
	wrappedHandler.ServeHTTP(rr, req)

	// Проверяем, что статус-код ответа 200 OK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Проверяем, что тело ответа содержит "OK"
	expected := "OK"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// Теперь тестируем случай, когда RateLimiter возвращает false
	handler.rate_limiter = &RateLimiter{
		Allow: func(hard int, dominantus bool) bool {
			return false
		},
	}

	// Создаем новый запрос и запись для ответа
	req = httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	rr = httptest.NewRecorder()

	// Вызываем обернутый обработчик
	wrappedHandler.ServeHTTP(rr, req)

	// Проверяем, что статус-код ответа 503 Service Unavailable
	if status := rr.Code; status != http.StatusServiceUnavailable {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusServiceUnavailable)
	}

	// Проверяем, что тело ответа содержит "Service Unavailable"
	expected = "Service Unavailable\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
