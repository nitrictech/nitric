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
	membrane, error := membrane.New()

	if error != nil {
		panic(error)
	}

	serviceAddress := getEnv("SERVICE_ADDRESS", "127.0.0.1:50051")
	childAddress := getEnv("CHILD_ADDRESS", "127.0.0.1:8080")
	childCommand := getEnv("INVOKE", "")

	// Start the Membrane server
	membrane.Start(serviceAddress, childAddress, childCommand)
}
