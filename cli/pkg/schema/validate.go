package schema

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/samber/lo"
	"github.com/xeipuuv/gojsonschema"
)

// Perform additional validation checks on the application
func (a *Application) IsValid() []gojsonschema.ResultError {
	// Check the names of all resources are unique
	violations := a.checkNoNameConflicts()
	violations = append(violations, a.checkNoReservedNames()...)
	violations = append(violations, a.checkSnakeCaseNames()...)
	violations = append(violations, a.checkNoEnvVarCollisions()...)
	violations = append(violations, a.checkAccessPermissions()...)

	return violations
}

func (a *Application) checkAccessPermissions() []gojsonschema.ResultError {
	violations := []gojsonschema.ResultError{}

	for name, intent := range a.BucketIntents {
		for serviceName, actions := range intent.Access {
			invalidActions, ok := hasInvalidActions(actions, []string{"read", "write", "delete"})
			if !ok {
				key := fmt.Sprintf("buckets.%s.access.%s", name, serviceName)
				err := fmt.Sprintf("provided invalid actions [%s]", strings.Join(invalidActions, ", "))
				violations = append(violations, newValidationError(key, err))
			}
		}
	}

	for name, intent := range a.DatabaseIntents {
		for serviceName, actions := range intent.Access {
			invalidActions, ok := hasInvalidActions(actions, []string{"query", "mutate"})
			if !ok {
				key := fmt.Sprintf("databases.%s.access.%s", name, serviceName)
				err := fmt.Sprintf("provided invalid actions [%s]", strings.Join(invalidActions, ", "))
				violations = append(violations, newValidationError(key, err))
			}
		}
	}

	return violations
}

func hasInvalidActions(actions []string, validActions []string) ([]string, bool) {
	invalidActions := []string{}
	validActions = append(validActions, "all")

	for _, action := range actions {
		if !lo.Contains(validActions, strings.ToLower(action)) {
			invalidActions = append(invalidActions, action)
		}
	}

	return invalidActions, len(invalidActions) == 0
}

func (a *Application) checkNoNameConflicts() []gojsonschema.ResultError {
	resourceNames := map[string]string{}
	violations := []gojsonschema.ResultError{}

	for name := range a.ServiceIntents {
		if existingType, ok := resourceNames[name]; ok {
			violations = append(violations, newValidationError(fmt.Sprintf("services.%s", name), fmt.Sprintf("service name %s is already in use by a %s", name, existingType)))
			continue
		}

		resourceNames[name] = "service"
	}

	for name := range a.BucketIntents {
		if existingType, ok := resourceNames[name]; ok {
			violations = append(violations, newValidationError(fmt.Sprintf("buckets.%s", name), fmt.Sprintf("bucket name %s is already in use by a %s", name, existingType)))
			continue
		}
		resourceNames[name] = "bucket"
	}

	for name := range a.EntrypointIntents {
		if existingType, ok := resourceNames[name]; ok {
			violations = append(violations, newValidationError(fmt.Sprintf("entrypoints.%s", name), fmt.Sprintf("entrypoint name %s is already in use by a %s", name, existingType)))
			continue
		}
		resourceNames[name] = "entrypoint"
	}

	for name := range a.DatabaseIntents {
		if existingType, ok := resourceNames[name]; ok {
			violations = append(violations, newValidationError(fmt.Sprintf("databases.%s", name), fmt.Sprintf("database name %s is already in use by a %s", name, existingType)))
			continue
		}
		resourceNames[name] = "database"
	}

	for name := range a.WebsiteIntents {
		if existingType, ok := resourceNames[name]; ok {
			violations = append(violations, newValidationError(fmt.Sprintf("websites.%s", name), fmt.Sprintf("website name %s is already in use by a %s", name, existingType)))
			continue
		}
		resourceNames[name] = "website"
	}

	return violations
}

func (a *Application) checkNoReservedNames() []gojsonschema.ResultError {
	violations := []gojsonschema.ResultError{}
	reservedNames := []string{
		"backend", // Backend is a reserved keyword in terraform
	}

	for name := range a.ServiceIntents {
		if slices.Contains(reservedNames, name) {
			violations = append(violations, newValidationError(fmt.Sprintf("services.%s", name), fmt.Sprintf("service name %s is a reserved name", name)))
		}
	}

	for name := range a.BucketIntents {
		if slices.Contains(reservedNames, name) {
			violations = append(violations, newValidationError(fmt.Sprintf("buckets.%s", name), fmt.Sprintf("bucket name %s is a reserved name", name)))
		}
	}

	for name := range a.EntrypointIntents {
		if slices.Contains(reservedNames, name) {
			violations = append(violations, newValidationError(fmt.Sprintf("entrypoints.%s", name), fmt.Sprintf("entrypoint name %s is a reserved name", name)))
		}
	}

	for name := range a.DatabaseIntents {
		if slices.Contains(reservedNames, name) {
			violations = append(violations, newValidationError(fmt.Sprintf("databases.%s", name), fmt.Sprintf("database name %s is a reserved name", name)))
		}
	}

	for name := range a.WebsiteIntents {
		if slices.Contains(reservedNames, name) {
			violations = append(violations, newValidationError(fmt.Sprintf("websites.%s", name), fmt.Sprintf("website name %s is a reserved name", name)))
		}
	}

	return violations
}

func (a *Application) checkSnakeCaseNames() []gojsonschema.ResultError {
	violations := []gojsonschema.ResultError{}
	snakeCasePattern := regexp.MustCompile(`^[a-z_][a-z0-9_]*$`)

	for name := range a.ServiceIntents {
		if !snakeCasePattern.MatchString(name) {
			violations = append(violations, newValidationError(fmt.Sprintf("services.%s", name), fmt.Sprintf("service name %s must be in snake_case format", name)))
		}
	}

	for name := range a.BucketIntents {
		if !snakeCasePattern.MatchString(name) {
			violations = append(violations, newValidationError(fmt.Sprintf("buckets.%s", name), fmt.Sprintf("bucket name %s must be in snake_case format", name)))
		}
	}

	for name := range a.EntrypointIntents {
		if !snakeCasePattern.MatchString(name) {
			violations = append(violations, newValidationError(fmt.Sprintf("entrypoints.%s", name), fmt.Sprintf("entrypoint name %s must be in snake_case format", name)))
		}
	}

	for name := range a.DatabaseIntents {
		if !snakeCasePattern.MatchString(name) {
			violations = append(violations, newValidationError(fmt.Sprintf("databases.%s", name), fmt.Sprintf("database name %s must be in snake_case format", name)))
		}
	}

	for name := range a.WebsiteIntents {
		if !snakeCasePattern.MatchString(name) {
			violations = append(violations, newValidationError(fmt.Sprintf("websites.%s", name), fmt.Sprintf("website name %s must be in snake_case format", name)))
		}
	}

	return violations
}

func (a *Application) checkNoEnvVarCollisions() []gojsonschema.ResultError {
	violations := []gojsonschema.ResultError{}
	envVarMap := map[string]string{}

	for name, intent := range a.DatabaseIntents {
		if existingName, ok := envVarMap[intent.EnvVarKey]; ok {
			violations = append(violations, newValidationError(fmt.Sprintf("databases.%s", name), fmt.Sprintf("env var %s is already in use by %s", intent.EnvVarKey, existingName)))
			continue
		}
		envVarMap[intent.EnvVarKey] = name
	}

	return violations
}
