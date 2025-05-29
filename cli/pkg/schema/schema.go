package schema

import (
	"encoding/json"
	"log"

	"github.com/invopop/jsonschema"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

type Application struct {
	// Targets sets platforms this application should be expected to work on
	// This gives us room to move away from LCD expectations around how platforms are built
	Targets     []string `json:"targets" yaml:"targets" jsonschema:"required"`
	Name        string   `json:"name" yaml:"name" jsonschema:"required"`
	Description string   `json:"description" yaml:"description"`

	Resources map[string]Resource `json:"resources,omitempty" yaml:"resources,omitempty"`
}

func ApplicationJsonSchema() *jsonschema.Schema {
	return jsonschema.Reflect(&Application{})
}

func ApplicationJsonSchemaString() string {
	schema := ApplicationJsonSchema()
	jsonSchemaOutput, err := json.Marshal(schema)

	if err != nil {
		log.Fatal(err.Error())
	}

	return string(jsonSchemaOutput)
}

func ApplicationFromYaml(yamlString string) (*Application, *gojsonschema.Result, error) {
	var raw map[string]interface{}

	err := yaml.Unmarshal([]byte(yamlString), &raw)
	if err != nil {
		return nil, nil, err
	}

	rawJson, err := json.Marshal(raw)

	return ApplicationFromJson(string(rawJson))
}

func ApplicationFromJson(jsonString string) (*Application, *gojsonschema.Result, error) {
	schemaLoader := gojsonschema.NewStringLoader(ApplicationJsonSchemaString())
	documentLoader := gojsonschema.NewStringLoader(jsonString)

	appSchema, _ := gojsonschema.NewSchema(schemaLoader)

	result, err := appSchema.Validate(documentLoader)
	if err != nil || !result.Valid() {
		return nil, result, err
	}

	var app Application
	err = json.Unmarshal([]byte(jsonString), &app)
	if err != nil {
		return nil, nil, err
	}

	return &app, nil, nil
}
