package schema

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

type NitricErrorTemplate struct {
	gojsonschema.DefaultLocale
}

type validationError struct {
	gojsonschema.ResultErrorFields
}

var _ gojsonschema.ResultError = &validationError{}

func newValidationError(field string, description string) gojsonschema.ResultError {
	e := &validationError{}

	e.SetDescription(description)
	e.SetType("invalid_property")
	e.SetContext(gojsonschema.NewJsonContext(field, nil))
	e.SetValue(map[string]interface{}{})

	return e
}

type ValidationError struct {
	Path        string `json:"path"`
	Message     string `json:"message"`
	ErrorType   string `json:"errorType"`
	YamlContext string `json:"yamlContext"`
}

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

	if path.ResourceType == "target" {
		return fmt.Sprintf("Invalid target at index %s", path.ResourceName)
	}

	if path.Property == "" {
		return fmt.Sprintf("%s %s has an invalid property", path.ResourceType, path.ResourceName)
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

func PrettyPrintPattern(pattern *regexp.Regexp) string {
	patterns := map[string]string{
		`^(([a-z]+)/([a-z]+)@(\d+)|file:([^\s]+))$`: "Must be in the format: `<team>/<platform>@<revision>` or `file:<path>`",
	}

	plainTextPattern, ok := patterns[pattern.String()]
	if !ok {
		return pattern.String()
	}

	return plainTextPattern
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

func (t *NitricErrorTemplate) DoesNotMatchPattern() string {
	return "{{pretty_print_pattern .pattern}}"
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
		"pretty_print_pattern": func(pattern *regexp.Regexp) string {
			return PrettyPrintPattern(pattern)
		},
	}
}

func GetSchemaValidationErrors(errs []gojsonschema.ResultError) []ValidationError {
	validationErrors := []ValidationError{}

	for _, err := range errs {
		formatter := newErrorFormatter(err.Field())

		// Skip certain errors based on context
		if formatter.ShouldSkipError(err.Type()) {
			continue
		}

		yamlContext := GenerateYamlContext(err, err.Description())

		validationErrors = append(validationErrors, ValidationError{
			Path:        err.Field(),
			Message:     err.String(),
			ErrorType:   err.Type(),
			YamlContext: yamlContext,
		})
	}

	slices.SortFunc(validationErrors, func(a, b ValidationError) int {
		return strings.Compare(a.Path, b.Path)
	})

	return validationErrors
}

/*
	Formats the error messages like this:

---> Error Message

	|
	| YAML Context
	|
*/
func FormatValidationErrors(errs []ValidationError) string {
	if len(errs) == 0 {
		return ""
	}

	arrow := "--->"
	linePrefix := "  |"
	var errsStr strings.Builder
	for _, err := range errs {
		errsStr.WriteString(fmt.Sprintf("%s %s\n", arrow, err.Message))

		for _, line := range strings.Split(err.YamlContext, "\n") {
			errsStr.WriteString(fmt.Sprintf("%s %s\n", linePrefix, line))
		}

		errsStr.WriteString(fmt.Sprintf("%s\n\n", linePrefix))
	}

	return errsStr.String()
}
