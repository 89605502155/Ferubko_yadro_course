package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	server "xkcd"
)

const (
	accessTimeConst  = time.Second
	refreshTimeConst = time.Hour
)

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func (h *Handler) SignIn(resp http.ResponseWriter, req *http.Request) {
	var data server.UserEntity

	// Прочитайте тело запроса
	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()

	// Декодируйте JSON в структуру
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(resp, "Bad Request", http.StatusBadRequest)
		return
	}
	err = data.Validate()
	if err != nil {
		http.Error(resp, "Unknow status", http.StatusBadRequest)
		return
	}
	access, refresh, err := h.services.Auth.GenerateToken(data, accessTimeConst, refreshTimeConst)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusUnauthorized)
		return
	}
	tokens := TokenResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}
	resp.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(resp).Encode(tokens)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	resp.WriteHeader(http.StatusOK)
}

func (h *Handler) CreateUser(resp http.ResponseWriter, req *http.Request) {
	logrus.Println("CreateUser", req.Method)
	var data server.User

	// Прочитайте тело запроса
	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()

	// Декодируйте JSON в структуру
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(resp, "Bad Request", http.StatusBadRequest)
		return
	}
	err = data.Validate()
	if err != nil {
		http.Error(resp, "Unknow status", http.StatusBadRequest)
		return
	}
	err = h.services.Auth.CreateUser(data)
	if err != nil {
		http.Error(resp, "Error in service", http.StatusBadRequest)
		return
	}
	// Используйте данные из структуры
	logrus.Printf("Received data: %+v\n", data)
	resp.Write([]byte("Good"))
}
