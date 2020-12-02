package main

import (
	"os"

	"nitric.io/membrane/membrane"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	serviceAddress := getEnv("SERVICE_ADDRESS", "127.0.0.1:50051")
	childAddress := getEnv("CHILD_ADDRESS", "127.0.0.1:8080")
	childCommand := getEnv("INVOKE", "echo No function configured")

	membrane, error := membrane.New(serviceAddress, childAddress, childCommand)

	if error != nil {
		panic(error)
	}

	// Start the Membrane server
	membrane.Start()
}
