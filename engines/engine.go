package engines

import "github.com/nitrictech/nitric/cli/pkg/schema"

// Common engine interface for all engines
type Engine interface {
	// Apply the engine to the target environment
	// TODO: Do we want the detailed platform schema to be owned by the engine?
	Apply(application *schema.Application, environment map[string]interface{}) error
}

type EngineConstructor = func(platform interface{}) Engine
