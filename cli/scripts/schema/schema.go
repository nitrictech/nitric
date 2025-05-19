package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/invopop/jsonschema"
	"github.com/nitrictech/nitric/cli/pkg/schema"
)

func main() {
	schema := jsonschema.Reflect(&schema.Application{})

	jsonOutput, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(string(jsonOutput))
}
