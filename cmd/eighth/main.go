package main

import (
	"log"

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
	defer database.Close()
	log.Println("db connected")

	srv, err := app.NewServer(&cfg.AppConfig)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("server is running on", cfg.AppConfig.Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
