package schema

import (
	"github.com/invopop/jsonschema"
)

type TargetType string

const (
	TargetType_Service TargetType = "service"
	TargetType_Website TargetType = "website"
)

type EntrypointIntent struct {
	Resource `json:",inline" yaml:",inline"`
	Routes   map[string]Route `json:"routes" yaml:"routes"`
}

func (e *EntrypointIntent) GetType() string {
	return "entrypoint"
}

func (e EntrypointIntent) JSONSchemaExtend(schema *jsonschema.Schema) {
	if routesSchema, ok := schema.Properties.Get("routes"); ok {
		routesSchema.PropertyNames = &jsonschema.Schema{
			Pattern: "/$",
		}
	}
}

type Route struct {
	TargetName string `json:"name" yaml:"name"`
	BasePath   string `json:"base-path,omitempty" yaml:"base-path,omitempty"`
}
