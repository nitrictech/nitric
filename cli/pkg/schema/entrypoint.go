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
	EntrypointSchemaOnlyHackType    string `json:"type" yaml:"-" jsonschema:"type,enum=entrypoint"`
	EntrypointSchemaOnlyHackSubType string `json:"sub-type,omitempty" yaml:"-,omitempty" jsonschema:"sub-type"`
	// TODO: As all resource names are unique, we could use the name as the value for the routes instead of the Route struct
	Routes map[string]Route `json:"routes" yaml:"routes"`
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
