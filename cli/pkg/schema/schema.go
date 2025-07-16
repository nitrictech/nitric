package schema

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/invopop/jsonschema"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

type Application struct {
	// Targets sets platforms this application should be expected to work on
	// This gives us room to move away from LCD expectations around how platforms are built
	Targets     []string `json:"targets" yaml:"targets" jsonschema:"required,pattern=^(([a-z]+)/([a-z]+)@(\\d+)|file:([^\\s]+))$"`
	Name        string   `json:"name" yaml:"name" jsonschema:"required"`
	Description string   `json:"description" yaml:"description"`

	ServiceIntents    map[string]*ServiceIntent    `json:"services,omitempty" yaml:"services,omitempty"`
	BucketIntents     map[string]*BucketIntent     `json:"buckets,omitempty" yaml:"buckets,omitempty"`
	EntrypointIntents map[string]*EntrypointIntent `json:"entrypoints,omitempty" yaml:"entrypoints,omitempty"`
	DatabaseIntents   map[string]*DatabaseIntent   `json:"databases,omitempty" yaml:"databases,omitempty"`
	WebsiteIntents    map[string]*WebsiteIntent    `json:"websites,omitempty" yaml:"websites,omitempty"`
}

func (a *Application) GetTypeForIntent(intent interface{}) (string, error) {
	switch intent.(type) {
	case *ServiceIntent:
		return "service", nil
	case *BucketIntent:
		return "bucket", nil
	case *EntrypointIntent:
		return "entrypoint", nil
	case *DatabaseIntent:
		return "database", nil
	case *WebsiteIntent:
		return "website", nil
	default:
		return "", fmt.Errorf("unknown intent type: %T", intent)
	}
}

func (a *Application) GetResourceIntents() map[string]IResource {
	resourceIntents := map[string]IResource{}

	for name, intent := range a.ServiceIntents {
		resourceIntents[name] = intent
	}

	for name, intent := range a.BucketIntents {
		resourceIntents[name] = intent
	}

	for name, intent := range a.EntrypointIntents {
		resourceIntents[name] = intent
	}

	for name, intent := range a.DatabaseIntents {
		resourceIntents[name] = intent
	}

	for name, intent := range a.WebsiteIntents {
		resourceIntents[name] = intent
	}

	return resourceIntents
}

func (a *Application) GetResourceIntent(name string) (interface{}, bool) {
	if service, ok := a.ServiceIntents[name]; ok {
		return service, true
	}

	if bucket, ok := a.BucketIntents[name]; ok {
		return bucket, true
	}

	if entrypoint, ok := a.EntrypointIntents[name]; ok {
		return entrypoint, true
	}

	if database, ok := a.DatabaseIntents[name]; ok {
		return database, true
	}

	if website, ok := a.WebsiteIntents[name]; ok {
		return website, true
	}

	return nil, false
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
	gojsonschema.ErrorTemplateFuncs = ErrorTemplateFunc
	gojsonschema.Locale = &NitricErrorTemplate{}

	schemaLoader := gojsonschema.NewStringLoader(ApplicationJsonSchemaString())
	documentLoader := gojsonschema.NewStringLoader(jsonString)

	var app Application
	err := json.Unmarshal([]byte(jsonString), &app)
	if err != nil {
		return nil, nil, err
	}

	appSchema, _ := gojsonschema.NewSchema(schemaLoader)

	result, err := appSchema.Validate(documentLoader)
	if err != nil || !result.Valid() {
		return &app, result, err
	}

	return &app, result, nil
}
