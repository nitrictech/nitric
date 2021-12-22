package deploy

import (
	"fmt"
	"strings"

	openapi "github.com/getkin/kin-openapi/openapi3"
)

type RouteHandler struct {
	path    string
	methods []string
}

type ApiBuilder struct {
	apis map[string][]*RouteHandler
}

func (ab *ApiBuilder) AddRouteHandler(api string, rh *RouteHandler) {
	rhs := ab.apis[api]

	if rhs == nil {
		rhs = make([]*RouteHandler, 0)
	}

	ab.apis[api] = append(rhs, rh)
}

// Translate an available API into an Open API v3 spec
func (ab *ApiBuilder) ToOaiSpec(api string) openapi.T {
	// Build open api specification from route definition
	doc := openapi.T{}
	defaultResponse := "Successful Response"
	doc.OpenAPI = "3.0.3"
	doc.Info = &openapi.Info{
		Title:   api,
		Version: "1",
	}

	rhs := ab.apis[api]

	for _, rh := range rhs {
		for _, m := range rh.methods {
			newOp := openapi.NewOperation()
			newOp.OperationID = fmt.Sprintf("%s-%s", rh.path, m)
			newOp.AddResponse(200, &openapi.Response{
				Description: &defaultResponse,
			})

			// Find and add parameters
			pathParts := strings.Split(rh.path, "/")
			// FIXME: Don't repeat this for every method...
			for _, p := range pathParts {
				if strings.HasPrefix(p, ":") {
					paramName := strings.Replace(p, ":", "", -1)
					newOp.AddParameter(&openapi.Parameter{
						In:       "path",
						Required: true,
						Name:     paramName,
					})
				}
			}

			doc.AddOperation(rh.path, m, newOp)
		}
	}

	return doc
}

func NewApiBuilder() *ApiBuilder {
	return &ApiBuilder{
		apis: make(map[string][]*RouteHandler),
	}
}
