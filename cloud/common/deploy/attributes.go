package deploy

import "fmt"

type CommonStackDetails struct {
	Project       string
	FullStackName string
	Stack         string
	Region        string
}

// Read nitric attributes from the provided deployment attributes
func CommonStackDetailsFromAttributes(attributes map[string]interface{}) (*CommonStackDetails, error) {
	iProject, hasProject := attributes["project"]
	project, isString := iProject.(string)
	if !hasProject || !isString || project == "" {
		// need a valid project name
		return nil, fmt.Errorf("project is not set or invalid")
	}

	iStack, hasStack := attributes["stack"]
	stack, isString := iStack.(string)
	if !hasStack || !isString || stack == "" {
		// need a valid stack name
		return nil, fmt.Errorf("stack is not set or invalid")
	}

	iRegion, hasRegion := attributes["region"]
	region, isString := iRegion.(string)
	if !hasRegion || !isString || stack == "" {
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
