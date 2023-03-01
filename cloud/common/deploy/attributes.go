package deploy

import "fmt"

type CommonStackDetails struct {
	Project       string
	FullStackName string
	Stack         string
	Region        string
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

	// Backwards compatible stack name
	// The existing providers in the CLI
	// Use the combined project and stack name
	fullStackName := fmt.Sprintf("%s-%s", project, stack)

	return &CommonStackDetails{
		Project:       project,
		FullStackName: fullStackName,
		Region:        region,
		Stack:         stack,
	}, nil
}
