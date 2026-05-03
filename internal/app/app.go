package app

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/BogdanBratsky/eigth/internal/config"
	"github.com/BogdanBratsky/eigth/internal/handler"
	"github.com/BogdanBratsky/eigth/internal/repository"
	"github.com/BogdanBratsky/eigth/internal/service"
	"github.com/BogdanBratsky/eigth/pkg/hasher"
	"github.com/BogdanBratsky/eigth/pkg/token"
)

func NewServer(appCfg *config.AppConfig, db *sql.DB) (*http.Server, error) {
	// infra
	hasher := hasher.NewBcryptHasher(10)
	jwt := token.NewJWTManager([]byte("he-he"), time.Minute)
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	loggerHandler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(loggerHandler)

	// repo
	userRepo := repository.NewUserPostgres(db)

	// service
	authService := service.NewAuthService(userRepo, hasher, jwt, logger)

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
