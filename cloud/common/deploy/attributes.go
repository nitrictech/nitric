package deploy

import "fmt"

type CommonStackDetails struct {
	Project string
	Stack   string
	Region  string
}

// Read nitric attributes from the provided deployment attributes
func CommonStackDetailsFromAttributes(attributes map[string]string) (*CommonStackDetails, error) {
	project, ok := attributes["project"]
	if !ok || project == "" {
		// need a valid project name
		return nil, fmt.Errorf("project is not set or invalid")
	}

	stack, ok := attributes["stack"]
	if !ok || stack == "" {
		// need a valid stack name
		return nil, fmt.Errorf("stack is not set or invalid")
	}

	region, ok := attributes["region"]
	if !ok || stack == "" {
		// need a valid stack name
		return nil, fmt.Errorf("region is not set or invalid")
	}

	return &CommonStackDetails{
		Project: project,
		Stack:   stack,
		Region:  region,
	}, nil
}
