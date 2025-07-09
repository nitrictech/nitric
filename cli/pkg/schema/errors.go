package schema

import (
	"fmt"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

type FieldPath struct {
	IsRoot       bool
	ResourceType string
	ResourceName string
	Property     string
	Value        string
}

func ParseFieldPath(field string) *FieldPath {
	splits := strings.Split(field, ".")

	if len(splits) < 2 {
		return &FieldPath{IsRoot: true}
	}

	path := &FieldPath{
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
		path.Value = splits[3]
	default:
		// For longer paths, take the last two elements
		if len(splits) > 3 {
			path.Property = splits[len(splits)-2]
			path.Value = splits[len(splits)-1]
		}
	}

	return path
}

type ErrorFormatter struct {
	path *FieldPath
}

func NewErrorFormatter(field string) *ErrorFormatter {
	return &ErrorFormatter{
		path: ParseFieldPath(field),
	}
}

func (ef *ErrorFormatter) FormatErrorPrefix() string {
	path := ef.path

	if path.IsRoot {
		return "Invalid application configuration"
	}

	if path.Property == "" {
		return fmt.Sprintf("%s %s has an invalid config", path.ResourceType, path.ResourceName)
	}

	if path.Value == "" {
		return fmt.Sprintf("%s %s has an invalid %s property", path.ResourceType, path.ResourceName, path.Property)
	}

	return fmt.Sprintf("%s %s has an invalid %s property (%s)", path.ResourceType, path.ResourceName, path.Property, path.Value)
}

func (ef *ErrorFormatter) FormatNumberOneOf() string {
	path := ef.path

	if path.ResourceType == "service" && path.Property == "container" {
		return "Must provide either a valid docker or image configuration. But not both."
	}

	return "Must validate one and only one schema"
}

func (ef *ErrorFormatter) FormatInvalidProperty() string {
	path := ef.path

	if path.ResourceType == "entrypoint" && path.Property == "routes" {
		return "Missing trailing slash for route"
	}

	return path.ResourceName
}

func (ef *ErrorFormatter) ShouldSkipError(errType string) bool {
	return errType == "pattern" && ef.path.ResourceType == "entrypoint" && ef.path.Property == "routes"
}

type NitricErrorTemplate struct {
	gojsonschema.DefaultLocale
}

func (t *NitricErrorTemplate) ErrorFormat() string {
	return "{{ error_prefix .field}}: {{.description}}"
}

func (t *NitricErrorTemplate) NumberOneOf() string {
	return "{{ one_of .field}}"
}

func (t *NitricErrorTemplate) InvalidPropertyName() string {
	return "{{ invalid_property_name .field}} {{.property}}"
}

func (t *NitricErrorTemplate) RegexPattern() string {
	return "{{ invalid_pattern .field}} {{.pattern}}"
}

func (t *NitricErrorTemplate) Required() string {
	return "The `{{.property}}` property is required"
}

func ErrorTemplateFunc() map[string]interface{} {
	return map[string]interface{}{
		"error_prefix": func(field string) string {
			formatter := NewErrorFormatter(field)
			return formatter.FormatErrorPrefix()
		},
		"one_of": func(field string) string {
			formatter := NewErrorFormatter(field)
			return formatter.FormatNumberOneOf()
		},
		"invalid_property_name": func(field string) string {
			formatter := NewErrorFormatter(field)
			return formatter.FormatInvalidProperty()
		},
	}
}

func FormatErrors(results *gojsonschema.Result) string {
	var errs strings.Builder

	for _, err := range results.Errors() {
		formatter := NewErrorFormatter(err.Field())

		// Skip certain errors based on context
		if formatter.ShouldSkipError(err.Type()) {
			continue
		}

		errs.WriteString(fmt.Sprintf(" - %s\n", err))
	}

	return errs.String()
}
