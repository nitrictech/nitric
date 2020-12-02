package main

import "nitric.io/membrane/membrane"

func main() {
	membrane, error := membrane.New()

	if error != nil {
		panic(error)
	}

	// Start the Membrane server
	membrane.Start()
}
