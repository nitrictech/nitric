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

	ResourceIntents map[string]Resource `json:"resources,omitempty" yaml:"resources,omitempty"`
}

func (a *Application) GetResourceIntentsForType(typ string) map[string]Resource {
	filteredResources := map[string]Resource{}

	for name, res := range a.ResourceIntents {
		if res.Type == typ {
			filteredResources[name] = res
		}
	}

	return filteredResources
}

func (a *Application) GetBucketIntents() map[string]*BucketIntent {
	concreteBuckets := map[string]*BucketIntent{}

	services := a.GetResourceIntentsForType("bucket")

	for name, svc := range services {
		concreteBuckets[name] = svc.BucketIntent
	}

	return concreteBuckets
}

func (a *Application) GetServiceIntents() map[string]*ServiceIntent {
	concreteServices := map[string]*ServiceIntent{}

	services := a.GetResourceIntentsForType("service")

	for name, svc := range services {
		concreteServices[name] = svc.ServiceIntent
	}

	return concreteServices
}

func (a *Application) GetEntrypointIntents() map[string]*EntrypointIntent {
	concreteEntrypoints := map[string]*EntrypointIntent{}

	services := a.GetResourceIntentsForType("entrypoint")

	for name, svc := range services {
		concreteEntrypoints[name] = svc.EntrypointIntent
	}

	return concreteEntrypoints
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
