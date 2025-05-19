package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"

	"github.com/invopop/jsonschema"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"

	"github.com/nitrictech/nitric/cli/pkg/schema"
)

//go:embed test.yaml
var testYaml []byte

//go:embed test.json
var testJson []byte

func main() {
	var appConfig schema.Application
	var rawConfig map[string]interface{}

	jsonSchema := jsonschema.Reflect(&schema.Application{})

	jsonSchemaOutput, err := json.MarshalIndent(jsonSchema, "", "  ")
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(string(jsonSchemaOutput))

	yaml.Unmarshal(testYaml, &appConfig)
	yaml.Unmarshal(testYaml, &rawConfig)

	actualRealJson, err := json.Marshal(&appConfig)
	if err != nil {
		panic(err)
	}

	if actualRealJson == nil {
		//
	}

	actualRealRawJson, err := json.Marshal(&rawConfig)
	if err != nil {
		panic(err)
	}

	// json.Unmarshal(testJson, &s)
	fmt.Printf("%+v\n", appConfig)

	for _, res := range appConfig.Resources {
		if res.Type == "service" {
			fmt.Printf("%+v\n", res.ServiceResource)
		} else if res.Type == "bucket" {
			fmt.Printf("%+v\n", res.BucketResource)
		}

	}

	schemaLoader := gojsonschema.NewStringLoader(string(jsonSchemaOutput))
	documentLoader := gojsonschema.NewStringLoader(string(actualRealRawJson))

	sc, _ := gojsonschema.NewSchema(schemaLoader)
	result, err := sc.Validate(documentLoader)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("valid: %v\n", result.Valid())
	fmt.Printf("errors: %+v\n", result.Errors())

	// compiler := boon.NewCompiler()
	// compiler.
	// 	data, _ := json.MarshalIndent(schema, "", "  ")

	// fmt.Printf("%s\n", string(data))
}
