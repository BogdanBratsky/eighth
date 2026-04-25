package app

import (
	"database/sql"
	"net/http"

	"github.com/BogdanBratsky/eigth/internal/config"
	"github.com/BogdanBratsky/eigth/internal/handler"
	"github.com/BogdanBratsky/eigth/internal/repository"
	"github.com/BogdanBratsky/eigth/internal/service"
	"github.com/BogdanBratsky/eigth/pkg/hasher"
)

func NewServer(appCfg *config.AppConfig, db *sql.DB) (*http.Server, error) {
	// repo
	userRepo := repository.NewUserPostgres(db)
	// infra
	hasher := hasher.NewBcryptHasher(10)
	// service
	authService := service.NewAuthService(userRepo, hasher)
	// handler
	authHandler := handler.NewAuthHandler(authService)

	// router
	mux := http.NewServeMux()
	// routes
	mux.HandleFunc("POST /register", authHandler.Register)

	return &http.Server{
		Addr:    appCfg.Port,
		Handler: mux,
	}, nil
}
