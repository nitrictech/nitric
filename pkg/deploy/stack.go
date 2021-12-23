package deploy

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	v1 "github.com/nitrictech/nitric/interfaces/nitric/v1"
)

// Stack - represents a collection of related functions and their shared dependencies.
type Stack struct {
	// A stack can be composed of one or more applications
	functions []*Function
}

// Produce an open api v3 spec for the requests API name
func (s *Stack) GetApiSpec(api string) (*openapi3.T, error) {
	doc := &openapi3.T{
		Paths: make(openapi3.Paths),
	}

	doc.Info = &openapi3.Info{
		Title:   api,
		Version: "v1",
	}

	doc.OpenAPI = "3.0.1"

	// Compile an API specification from the functions in the stack for the given API name
	workers := make([]*v1.ApiWorker, 0)

	// Collect all workers
	for _, f := range s.functions {

		workers = append(workers, f.apis[api].workers...)
	}

	// loop over workers to build new api specification
	// FIXME: We will need to merge path matches across all workers
	// to ensure we don't have conflicts
	for _, w := range workers {
		params := make(openapi3.Parameters, 0)
		normalizedPath := ""
		for _, p := range strings.Split(w.Path, "/") {
			if strings.HasPrefix(p, ":") {
				paramName := strings.Replace(p, ":", "", -1)
				params = append(params, &openapi3.ParameterRef{
					Value: &openapi3.Parameter{
						In:   "path",
						Name: paramName,
					},
				})
				normalizedPath = normalizedPath + "{" + paramName + "}" + "/"
			} else {
				normalizedPath = normalizedPath + p + "/"
			}
		}

		pathItem := doc.Paths.Find(normalizedPath)

		if pathItem == nil {
			// Add the parameters at the path level
			pathItem = &openapi3.PathItem{
				Parameters: params,
			}
			// Add the path item to the document
			doc.Paths[normalizedPath] = pathItem
		}

		for _, m := range w.Methods {
			if pathItem.Operations() != nil && pathItem.Operations()[m] != nil {
				// If the operation already exists we should fail
				// NOTE: This should not happen as operations are stored in a map
				// in the api state for functions
				return nil, fmt.Errorf("found conflicting operations")
			}

			// See if the path already exists
			doc.AddOperation(normalizedPath, m, &openapi3.Operation{
				OperationID: normalizedPath + m,
				Responses:   openapi3.NewResponses(),
			})
		}
	}

	return doc, nil
}

// Listen - Starts server to listen for new resource registrations for the stack
func (s *Stack) Listen() {

}
