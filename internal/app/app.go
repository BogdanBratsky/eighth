package app

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/BogdanBratsky/eigth/internal/config"
	"github.com/BogdanBratsky/eigth/internal/handler"
	"github.com/BogdanBratsky/eigth/internal/repository"
	"github.com/BogdanBratsky/eigth/internal/service"
	"github.com/BogdanBratsky/eigth/pkg/hasher"
	"github.com/BogdanBratsky/eigth/pkg/token"
)

func NewServer(appCfg *config.AppConfig, db *sql.DB) (*http.Server, error) {
	// repo
	userRepo := repository.NewUserPostgres(db)
	// infra
	hasher := hasher.NewBcryptHasher(10)
	jwt := token.NewJWTManager([]byte("he-he"), time.Minute)
	// service
	authService := service.NewAuthService(userRepo, hasher, jwt)
	// handler
	authHandler := handler.NewAuthHandler(authService)

	// router
	mux := http.NewServeMux()
	// routes
	mux.HandleFunc("POST /register", authHandler.Register)
	mux.HandleFunc("POST /login", authHandler.Login)

	return &http.Server{
		Addr:    appCfg.Port,
		Handler: mux,
	}, nil
}
