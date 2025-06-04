package schema

import "github.com/invopop/jsonschema"

type Resource struct {
	Type    string `json:"type" yaml:"type" jsonschema:"-"`
	SubType string `json:"sub-type,omitempty" yaml:"sub-type,omitempty"`

	// A resource can contain oneof the following sets of keys (see JSONSchemaExtended)
	*ServiceIntent      `json:",inline,omitempty" yaml:",inline,omitempty" jsonschema:"-"`
	*BucketIntent       `json:",inline,omitempty" yaml:",inline,omitempty" jsonschema:"-"`
	*EntrypointIntent   `json:",inline,omitempty" yaml:",inline,omitempty" jsonschema:"-"`
	*SubscriptionIntent `json:",inline,omitempty" yaml:",inline,omitempty" jsonschema:"-"`
	*DatabaseIntent     `json:",inline,omitempty" yaml:",inline,omitempty" jsonschema:"-"`
	*StateIntent        `json:",inline,omitempty" yaml:",inline,omitempty" jsonschema:"-"`
}

// schema types defined for the output schema
var schemaTypes = map[string]interface{}{
	"ServiceResource":      ServiceIntent{},
	"BucketResource":       BucketIntent{},
	"EntrypointResource":   EntrypointIntent{},
	"SubscriptionResource": SubscriptionIntent{},
	"DatabaseResource":     DatabaseIntent{},
	"StateResource":        StateIntent{},
}

func (Resource) JSONSchemaExtend(schema *jsonschema.Schema) {
	if schema.Definitions == nil {
		schema.Definitions = map[string]*jsonschema.Schema{}
	}

	subSchemas := []*jsonschema.Schema{}
	for _, res := range schemaTypes {
		s := jsonschema.Reflect(res)

		s.AdditionalProperties = nil
		s.Properties = nil

		subSchemas = append(subSchemas, s)
	}

	schema.Properties = nil
	schema.AdditionalProperties = nil
	schema.OneOf = subSchemas
}
