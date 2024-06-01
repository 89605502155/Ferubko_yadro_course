package handler

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"

	server "xkcd"
)

func (h *Handler) SignIn(resp http.ResponseWriter, req *http.Request) {
	logrus.Println("Auth", req.Method)
	resp.Write([]byte("You are user\n"))
	resp.Header().Set("Token", "<PASSWORD>")
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
