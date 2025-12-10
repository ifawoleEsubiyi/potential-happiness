// Package main provides the entry point for the validatord command.
// This is a thin wrapper that invokes the main validatord daemon functionality.
package main

import (
	"log"

	"github.com/dreadwitdastacc-IFA/validatord/internal/app"
)

func main() {
	log.Println("Starting validatord...")

	application, err := app.New(app.DefaultPaystring)
	if err != nil {
		log.Fatalf("Failed to initialize validatord: %v", err)
	}

	application.PrintStatus()
}
