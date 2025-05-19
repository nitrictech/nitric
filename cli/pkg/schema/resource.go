package schema

import "github.com/invopop/jsonschema"

type Resource struct {
	Type    string `json:"type" yaml:"type" jsonschema:"-"`
	SubType string `json:"sub-type,omitempty" yaml:"sub-type,omitempty"`

	// A resource can contain oneof the following sets of keys (see JSONSchemaExtended)
	*ServiceResource      `json:",inline,omitempty" yaml:",inline,omitempty" jsonschema:"-"`
	*BucketResource       `json:",inline,omitempty" yaml:",inline,omitempty" jsonschema:"-"`
	*EntrypointResource   `json:",inline,omitempty" yaml:",inline,omitempty" jsonschema:"-"`
	*SubscriptionResource `json:",inline,omitempty" yaml:",inline,omitempty" jsonschema:"-"`
	*DatabaseResource `json:",inline,omitempty" yaml:",inline,omitempty" jsonschema:"-"`
	*StateResource    `json:",inline,omitempty" yaml:",inline,omitempty" jsonschema:"-"`
}

// schema types defined for the output schema
var schemaTypes = map[string]interface{}{
	"ServiceResource":      ServiceResource{},
	"BucketResource":       BucketResource{},
	"EntrypointResource":   EntrypointResource{},
	"SubscriptionResource": SubscriptionResource{},
	"DatabaseResource": DatabaseResource{},
	"StateResource":    StateResource{},
}

func (Resource) JSONSchemaExtend(schema *jsonschema.Schema) {
	if schema.Definitions == nil {
		schema.Definitions = map[string]*jsonschema.Schema{}
	}

	subSchemas := []*jsonschema.Schema{}
	for _, res := range resourcesTypes {
		s := jsonschema.Reflect(res)

		s.AdditionalProperties = nil
		s.Properties = nil

		subSchemas = append(subSchemas, s)
	}

	schema.Properties = nil
	schema.AdditionalProperties = nil
	schema.OneOf = subSchemas
}
