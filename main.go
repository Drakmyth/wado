package main

import (
	"log"
	"log/slog"

	"github.com/Drakmyth/go-template/server"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Info("No .env file found. Using environment or defaults...")
	} else {
		slog.Info("Found .env file.")
	}

	log.Fatal(server.ListenAndServe())
}
