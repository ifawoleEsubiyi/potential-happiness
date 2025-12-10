// Package main is the entry point for the validatord daemon.
// Validatord is a validator daemon that provides infrastructure for blockchain validators,
// including attestation, BLS signature operations, key management, data aggregation, and state monitoring.
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
