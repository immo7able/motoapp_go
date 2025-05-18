package main

import (
	"log"
	"motorcycleApp/config"
	"motorcycleApp/server"
)

func main() {

	app := server.SetupApp()

	cfg := config.NewConfig()
	if err := app.Run(cfg.ServerAddress); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
