package schema

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/nitrictech/nitric/cli/internal/style/colors"
	"github.com/xeipuuv/gojsonschema"
)

type YamlContextBuilder struct {
	strings.Builder
}

func (b *YamlContextBuilder) WriteError(s string) (int, error) {
	errStr := lipgloss.NewStyle().Foreground(colors.Red).Render(fmt.Sprintf("\t# <-- %s", s))
	return b.Builder.WriteString(errStr)
}

func (b *YamlContextBuilder) WriteYamlKey(s string) (int, error) {
	keyStr := lipgloss.NewStyle().Foreground(colors.Blue).Render(s)
	return b.Builder.WriteString(keyStr)
}

func (b *YamlContextBuilder) String() string {
	return strings.TrimSpace(b.Builder.String())
}

// GenerateYamlContext generates the YAML context for validation errors
func GenerateYamlContext(err gojsonschema.ResultError) string {
	path := err.Field()
	yaml := YamlContextBuilder{}

	parts := strings.Split(path, ".")
	if len(parts) > 0 {
		for i, part := range parts {
			yaml.WriteYamlKey(fmt.Sprintf("%s%s:", strings.Repeat("  ", i), part))

			if i == len(parts)-1 {
				// Add offending property key
				if err.Type() == "required" || err.Type() == "additional_property_not_allowed" {
					yaml.WriteYamlKey(fmt.Sprintf("\n%s%s:", strings.Repeat("  ", i+1), err.Details()["property"]))
				}

				if errMap, ok := err.Value().(map[string]interface{}); ok {
					if len(errMap) > 0 {
						addErrorContext(&yaml, err, i)
						writeMapValue(&yaml, errMap, i)
					} else {
						addErrorContext(&yaml, err, i)
					}
				} else {
					yaml.WriteString(fmt.Sprintf(" %v", err.Value()))
					addErrorContext(&yaml, err, i)
				}
			} else {
				yaml.WriteString("\n")
			}
		}
	}

	return yaml.String()
}

func writeValue(contextBuilder *YamlContextBuilder, err gojsonschema.ResultError, indent int) {

}

func writeMapValue(contextBuilder *YamlContextBuilder, errMap map[string]interface{}, indent int) {
	// Only print if the value is not a map (i.e. a primitive type)
	if len(errMap) > 0 {
		for key, value := range errMap {
			contextBuilder.WriteYamlKey(fmt.Sprintf("\n%s%s:", strings.Repeat("  ", indent+1), key))

			if _, ok := value.(map[string]interface{}); ok {
				writeMapValue(contextBuilder, value.(map[string]interface{}), indent+1)
			} else {
				contextBuilder.WriteString(fmt.Sprintf(" %v", value))
			}
		}
	} else {
		contextBuilder.WriteString(" {}")
	}
}

func addErrorContext(contextBuilder *YamlContextBuilder, err gojsonschema.ResultError, indent int) {
	details := err.Details()

	switch err.Type() {
	case "required":
		if property, ok := details["property"].(string); ok {
			contextBuilder.WriteError(fmt.Sprintf("Missing %s", property))
		} else {
			contextBuilder.WriteError("Missing required property")
		}

	case "invalid_type":
		if expected, ok := details["expected"].(string); ok {
			contextBuilder.WriteError(fmt.Sprintf("Must be %s", expected))
		} else {
			contextBuilder.WriteError("Invalid type")
		}

	case "enum":
		if allowed, ok := details["allowed"].([]interface{}); ok {
			values := make([]string, len(allowed))
			for i, v := range allowed {
				if str, ok := v.(string); ok {
					values[i] = str
				}
			}
			contextBuilder.WriteError(fmt.Sprintf("Must be one of: %s", strings.Join(values, ", ")))
		} else {
			contextBuilder.WriteError("Invalid enum value")
		}

	case "number_one_of":
		contextBuilder.WriteError(fmt.Sprintf("Configuration must match exactly one schema."))

	case "number_any_of":
		contextBuilder.WriteError("Configuration must match at least one schema.")

	case "number_all_of":
		contextBuilder.WriteError("Configuration must match all schemas.")

	case "const":
		contextBuilder.WriteError("Value must be a specific constant")

	case "array_min_items":
		contextBuilder.WriteError("Array has too few items")

	case "array_max_items":
		contextBuilder.WriteError("Array has too many items")

	case "unique":
		contextBuilder.WriteError("Array items must be unique")

	case "additional_property_not_allowed":
		contextBuilder.WriteError("Additional properties not allowed")

	case "pattern":
		contextBuilder.WriteError("String does not match required pattern")

	case "string_gte":
		if min, ok := details["min"].(float64); ok {
			contextBuilder.WriteError(fmt.Sprintf("Must be >= %.0f characters", min))
		} else {
			contextBuilder.WriteError("String too short")
		}

	case "string_lte":
		if max, ok := details["max"].(float64); ok {
			contextBuilder.WriteError(fmt.Sprintf("Must be <= %.0f characters", max))
		} else {
			contextBuilder.WriteError("String too long")
		}

	case "number_gte":
		if min, ok := details["min"].(float64); ok {
			contextBuilder.WriteError(fmt.Sprintf("Must be >= %.0f", min))
		} else {
			contextBuilder.WriteError("Number too small")
		}

	case "number_gt":
		if min, ok := details["min"].(float64); ok {
			contextBuilder.WriteError(fmt.Sprintf("Must be > %.0f", min))
		} else {
			contextBuilder.WriteError("Number too small")
		}

	case "number_lte":
		if max, ok := details["max"].(float64); ok {
			contextBuilder.WriteError(fmt.Sprintf("Must be <= %.0f", max))
		} else {
			contextBuilder.WriteError("Number too large")
		}

	case "number_lt":
		if max, ok := details["max"].(float64); ok {
			contextBuilder.WriteError(fmt.Sprintf("Must be < %.0f", max))
		} else {
			contextBuilder.WriteError("Number too large")
		}

	case "multiple_of":
		if multipleOf, ok := details["multiple_of"].(float64); ok {
			contextBuilder.WriteError(fmt.Sprintf("Must be multiple of %.0f", multipleOf))
		} else {
			contextBuilder.WriteError("Number must be a multiple of specified value")
		}

	default:
		contextBuilder.WriteError(fmt.Sprintf("%s", strings.ReplaceAll(err.Type(), "_", " ")))
	}
}
