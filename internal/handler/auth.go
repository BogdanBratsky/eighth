package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/BogdanBratsky/eigth/internal/service"
)

type AuthHandler struct {
	service *service.AuthService
	logger  *slog.Logger
}

func NewAuthHandler(s *service.AuthService, l *slog.Logger) *AuthHandler {
	return &AuthHandler{
		service: s,
		logger:  l,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	h.logger.Info("register request received",
		"method", r.Method,
		"path", r.URL.Path,
	)

	var input RegisterReq
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.logger.Error("failed to decode register request",
			"error", err,
		)
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	err := h.service.Register(ctx, input.Login, input.Email, input.Password)
	if err != nil {
		h.logger.Error("register failed",
			"email", input.Email,
			"error", err,
		)
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.logger.Info("user registered successfully",
		"email", input.Email,
	)

	respondSuccess(w, http.StatusOK, nil)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	h.logger.Info("login request received",
		"method", r.Method,
		"path", r.URL.Path,
	)

	var input LoginReq
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.logger.Error("failed to decode login request",
			"error", err,
		)
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.service.Login(ctx, input.Email, input.Password)
	if err != nil {
		h.logger.Error("login failed",
			"email", input.Email,
			"error", err,
		)
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	h.logger.Info("login successful",
		"email", input.Email,
	)

	respondSuccess(w, http.StatusOK, token)
}
