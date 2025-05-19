package schema

type TargetType string

const (
	TargetType_Service TargetType = "service"
	TargetType_Website TargetType = "website"
)

type EntrypointResource struct {
	EntrypointSchemaOnlyHackType string `json:"type" jsonschema:"type,enum=entrypoint"`
	// TODO: As all resource names are unique, we could use the name as the value for the routes instead of the Route struct
	Routes map[string]Route `json:"routes"`
}

type Route struct {
	TargetName string `json:"name"`
}
