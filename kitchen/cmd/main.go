package main

import (
	"context"
	"log"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/pkg/app"
)

func main() {
	application, err := app.New(context.Background(), app.ConfigFromEnv())
	if err != nil {
		log.Fatal(err)
	}
	defer application.Close()

	if err := application.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
