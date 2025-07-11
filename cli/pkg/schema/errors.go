package schema

import (
	"fmt"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

type fieldPath struct {
	IsRoot       bool
	ResourceType string
	ResourceName string
	Property     string
	SubProperty  string
}

func parseFieldPath(field string) *fieldPath {
	splits := strings.Split(field, ".")

	if len(splits) < 2 {
		return &fieldPath{IsRoot: true}
	}

	path := &fieldPath{
		ResourceType: strings.TrimSuffix(splits[0], "s"),
		ResourceName: splits[1],
	}

	switch len(splits) {
	case 2:
		return path
	case 3:
		path.Property = splits[2]
	case 4:
		path.Property = splits[2]
		path.SubProperty = splits[3]
	default:
		// For longer paths, take the last two elements
		if len(splits) > 3 {
			path.Property = splits[len(splits)-2]
			path.SubProperty = splits[len(splits)-1]
		}
	}

	return path
}

type errorFormatter struct {
	path *fieldPath
}

func newErrorFormatter(field string) *errorFormatter {
	return &errorFormatter{
		path: parseFieldPath(field),
	}
}

func (ef *errorFormatter) FormatErrorPrefix() string {
	path := ef.path

	if path.IsRoot {
		return "Invalid application configuration"
	}

	if path.Property == "" {
		return fmt.Sprintf("%s %s has an invalid config", path.ResourceType, path.ResourceName)
	}

	if path.SubProperty == "" {
		return fmt.Sprintf("%s %s has an invalid %s property", path.ResourceType, path.ResourceName, path.Property)
	}

	return fmt.Sprintf("%s %s has an invalid %s property (%s)", path.ResourceType, path.ResourceName, path.Property, path.SubProperty)
}

func (ef *errorFormatter) FormatInvalidProperty() string {
	path := ef.path

	if path.ResourceType == "entrypoint" && path.Property == "routes" {
		return "Missing trailing slash for route"
	}

	return path.ResourceName
}

func (ef *errorFormatter) FormatNumberOneOf() string {
	path := ef.path
	schemas := []string{}

	if path.ResourceType == "service" && path.Property == "container" {
		schemas = append(schemas, "docker", "image")
	}

	return fmt.Sprintf("Must provide exactly one of: %s", strings.Join(schemas, " OR "))
}

func (ef *errorFormatter) ShouldSkipError(errType string) bool {
	return errType == "pattern" && ef.path.ResourceType == "entrypoint" && ef.path.Property == "routes"
}

type NitricErrorTemplate struct {
	gojsonschema.DefaultLocale
}

func (t *NitricErrorTemplate) ErrorFormat() string {
	return "{{ error_prefix .field}}: {{.description}}"
}

func (t *NitricErrorTemplate) NumberOneOf() string {
	return "{{one_of .field}}"
}

func (t *NitricErrorTemplate) InvalidPropertyName() string {
	return "{{ invalid_property_name .field}} {{.property}}"
}

func (t *NitricErrorTemplate) RegexPattern() string {
	return "{{ invalid_pattern .field}} {{.pattern}}"
}

func (t *NitricErrorTemplate) Required() string {
	return "The {{.property}} property is required"
}

func ErrorTemplateFunc() map[string]interface{} {
	return map[string]interface{}{
		"error_prefix": func(field string) string {
			formatter := newErrorFormatter(field)
			return formatter.FormatErrorPrefix()
		},
		"invalid_property_name": func(field string) string {
			formatter := newErrorFormatter(field)
			return formatter.FormatInvalidProperty()
		},
		"one_of": func(field string) string {
			formatter := newErrorFormatter(field)
			return formatter.FormatNumberOneOf()
		},
	}
}

type ValidationError struct {
	Path        string `json:"path"`
	Message     string `json:"message"`
	ErrorType   string `json:"errorType"`
	YamlContext string `json:"yamlContext"`
}

func GetSchemaValidationErrors(results *gojsonschema.Result) []ValidationError {
	errs := []ValidationError{}

	for _, err := range results.Errors() {
		formatter := newErrorFormatter(err.Field())

		// Skip certain errors based on context
		if formatter.ShouldSkipError(err.Type()) {
			continue
		}

		yamlContext := GenerateYamlContext(err)

		errs = append(errs, ValidationError{
			Path:        err.Field(),
			Message:     err.String(),
			ErrorType:   err.Type(),
			YamlContext: yamlContext,
		})
	}

	return errs
}

/*
	Formats the error messages like this:

---> Error Message

	|
	| YAML Context
	|
*/
func FormatValidationErrors(results *gojsonschema.Result) string {
	errs := GetSchemaValidationErrors(results)

	arrow := "--->"
	linePrefix := "  |"
	var errsStr strings.Builder
	for _, err := range errs {
		lines := strings.Split(err.YamlContext, "\n")
		errsStr.WriteString(fmt.Sprintf("%s %s\n%s\n", arrow, err.Message, linePrefix))
		for _, line := range lines {
			errsStr.WriteString(fmt.Sprintf("%s %s\n", linePrefix, line))
		}
		errsStr.WriteString(fmt.Sprintf("%s\n\n", linePrefix))
	}

	return errsStr.String()
}
