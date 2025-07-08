package schema

import (
	"fmt"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

type FieldInformation struct {
	Type          string
	Name          string
	Property      string
	OptionalValue string
}

func getInformationFromField(field string) *FieldInformation {
	splits := strings.Split(field, ".")
	if len(splits) < 2 {
		return nil
	}

	fieldInfo := &FieldInformation{}

	fieldInfo.Type = strings.TrimSuffix(splits[0], "s")
	fieldInfo.Name = splits[1]

	if len(splits) > 2 {
		fieldInfo.Property = splits[len(splits)-1]
	}

	if len(splits) > 4 {
		fieldInfo.OptionalValue = splits[len(splits)-1]
	}

	return fieldInfo
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

func ErrorTemplateFunc() map[string]interface{} {
	return map[string]interface{}{
		"error_prefix": func(field string) string {
			resource := getInformationFromField(field)

			if resource.Type == "(root)" {
				return fmt.Sprintf("Invalid application config")
			}

			if resource.Property == "" {
				return fmt.Sprintf("%s %s has invalid config", resource.Type, resource.Name)
			}

			if resource.OptionalValue == "" {
				return fmt.Sprintf("%s %s has invalid %s property", resource.Type, resource.Name, resource.Property)
			}

			return fmt.Sprintf("%s %s has invalid %s property (%s)", resource.Type, resource.Name, resource.Property, resource.OptionalValue)
		},
		"one_of": func(field string) string {
			resource := getInformationFromField(field)

			if resource.Type == "service" && resource.Property == "container" {
				return "Must provide a valid docker or image configuration, but not both."
			}

			return "Must validate one and only one schema"
		},
		"invalid_property_name": func(field string) string {
			resource := getInformationFromField(field)

			if resource.Type == "entrypoint" && resource.Property == "route" {
				return "Missing trailing slash for route:"
			}

			return field
		},
	}
}

func FormatErrors(results *gojsonschema.Result) string {
	errs := ""
	for _, err := range results.Errors() {
		resource := getInformationFromField(err.Field())

		// Ignore printing the pattern matching error for entrypoint routes (handled by the invalid_property_name function)
		if err.Type() == "pattern" && resource.Type == "entrypoint" && resource.Property == "route" {
			continue
		}

		errs += fmt.Sprintf(" - %s\n", err)
	}

	return errs
}
