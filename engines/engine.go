package engines

import "github.com/nitrictech/nitric/cli/pkg/schema"

// Common engine interface for all engines
type Engine interface {
	Apply(application *schema.Application) error
}

type EngineConstructor = func(platform interface{}) Engine
