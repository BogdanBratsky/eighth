package app

import (
	"net/http"

	"github.com/BogdanBratsky/eigth/internal/config"
)

func NewServer(appCfg *config.AppConfig) (*http.Server, error) {
	mux := http.NewServeMux()
	return &http.Server{
		Addr:    appCfg.Port,
		Handler: mux,
	}, nil
}
