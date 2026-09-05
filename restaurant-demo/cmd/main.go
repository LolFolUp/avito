package main

import (
	"context"
	"log"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/pkg/app"
)

func main() {
	config, err := app.ConfigFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	application, err := app.New(config)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := application.Close(); err != nil {
			log.Printf("close application: %v", err)
		}
	}()

	if err := application.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
