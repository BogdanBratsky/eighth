package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BogdanBratsky/eigth/internal/app"
	"github.com/BogdanBratsky/eigth/internal/config"
	"github.com/BogdanBratsky/eigth/internal/db"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	log.Println("loading config")
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("config was load")

	log.Println("connecting to db")
	database, err := db.Connect(&cfg.DBConfig)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("db connected")

	srv, err := app.NewServer(&cfg.AppConfig, database)
	if err != nil {
		log.Fatal(err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("server is running on", cfg.AppConfig.Port)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	<-stop
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Println("server shutdown error:", err)
	}

	log.Println("closing db...")
	if err := database.Close(); err != nil {
		log.Println("db close error:", err)
	}

	log.Println("server stopped gracefully")
}
