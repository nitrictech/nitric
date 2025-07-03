package schema

import (
	"encoding/json"
	"errors"
	"fmt"
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

	ServiceIntents    map[string]*ServiceIntent    `json:"services,omitempty" yaml:"services,omitempty"`
	BucketIntents     map[string]*BucketIntent     `json:"buckets,omitempty" yaml:"buckets,omitempty"`
	EntrypointIntents map[string]*EntrypointIntent `json:"entrypoints,omitempty" yaml:"entrypoints,omitempty"`
	DatabaseIntents   map[string]*DatabaseIntent   `json:"databases,omitempty" yaml:"databases,omitempty"`
	WebsiteIntents    map[string]*WebsiteIntent    `json:"websites,omitempty" yaml:"websites,omitempty"`
}

// Perform additional validation checks on the application
func (a *Application) IsValid() error {
	// Check the names of all resources are unique
	resourceNames := map[string]string{}
	violations := []error{}

	for name, _ := range a.ServiceIntents {
		if existingType, ok := resourceNames[name]; ok {
			violations = append(violations, fmt.Errorf("service name %s is already in use by a %s", name, existingType))
			continue
		}
		resourceNames[name] = "service"
	}

	for name, _ := range a.BucketIntents {
		if existingType, ok := resourceNames[name]; ok {
			violations = append(violations, fmt.Errorf("bucket name %s is already in use by a %s", name, existingType))
			continue
		}
		resourceNames[name] = "bucket"
	}

	for name, _ := range a.EntrypointIntents {
		if existingType, ok := resourceNames[name]; ok {
			violations = append(violations, fmt.Errorf("entrypoint name %s is already in use by a %s", name, existingType))
			continue
		}
		resourceNames[name] = "entrypoint"
	}

	for name, _ := range a.DatabaseIntents {
		if existingType, ok := resourceNames[name]; ok {
			violations = append(violations, fmt.Errorf("database name %s is already in use by a %s", name, existingType))
			continue
		}
		resourceNames[name] = "database"
	}

	for name, _ := range a.WebsiteIntents {
		if existingType, ok := resourceNames[name]; ok {
			violations = append(violations, fmt.Errorf("website name %s is already in use by a %s", name, existingType))
			continue
		}
		resourceNames[name] = "website"
	}

	if len(violations) > 0 {
		// format the violations as a list
		violationsString := "Errors found in application:\n"
		for _, violation := range violations {
			violationsString += fmt.Sprintf(" - %s\n", violation)
		}

		return errors.New(violationsString)
	}

	return nil
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
	schemaLoader := gojsonschema.NewStringLoader(ApplicationJsonSchemaString())
	documentLoader := gojsonschema.NewStringLoader(jsonString)

	appSchema, _ := gojsonschema.NewSchema(schemaLoader)

	var app Application
	err := json.Unmarshal([]byte(jsonString), &app)
	if err != nil {
		return nil, nil, err
	}

	result, err := appSchema.Validate(documentLoader)
	if err != nil || !result.Valid() {
		return &app, result, err
	}

	return &app, nil, nil
}
