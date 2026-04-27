package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

type response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type errorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func success(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(response{
		Status:  status,
		Message: "success",
		Data:    data,
	})
	if err != nil {
		log.Println(err)
		return
	}
}
