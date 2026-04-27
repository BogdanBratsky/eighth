package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/BogdanBratsky/eigth/internal/service"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	log.Println("request recieved:", r.Method, r.URL)

	ctx := r.Context()

	var input RegisterReq
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.service.Register(ctx, input.Login, input.Email, input.Password)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Println("success")
	respondSuccess(w, http.StatusOK, nil)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	log.Println("request recieved:", r.Method, r.URL)

	ctx := r.Context()

	var input LoginReq
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	token, err := h.service.Login(ctx, input.Email, input.Password)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Println("success")
	respondSuccess(w, http.StatusOK, token)
}
