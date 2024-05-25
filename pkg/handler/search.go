package handler

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

func (h *Handler) Search(resp http.ResponseWriter, req *http.Request) {
	s := req.URL.Query().Get("search")
	rers1, t1, err := h.services.Search.SearchInDB(s)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}
	rers2, t2, err := h.services.Search.SearchInIndex(s)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}
	resp.Header().Add("Content-Type", "application/json")

	res := make(map[string]interface{})
	res["in db"] = rers1
	res["in index"] = rers2
	jn, _ := json.Marshal(res)
	resp.Write(jn)

	logrus.Println(t1, t2)
	resp.Header().Set("Content-Type", "application/json")
}
