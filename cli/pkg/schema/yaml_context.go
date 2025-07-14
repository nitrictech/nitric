package schema

import (
	"fmt"
	"strconv"
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

func (b *YamlContextBuilder) WriteYamlKey(s interface{}, indent int) (int, error) {
	indentStr := strings.Repeat("  ", indent)
	// If the key is an integer, it means its an array index, replace it with a dash.
	keyStr := ""
	if _, err := strconv.Atoi(fmt.Sprintf("%v", s)); err == nil {
		keyStr = fmt.Sprintf("\n%s-", indentStr)
	} else {
		keyStr = fmt.Sprintf("\n%s%s:", indentStr, s)
	}

	styledKeyStr := lipgloss.NewStyle().Foreground(colors.Blue).Render(keyStr)
	return b.Builder.WriteString(styledKeyStr)
}

func (b *YamlContextBuilder) String() string {
	return strings.TrimSpace(b.Builder.String())
}

// GenerateYamlContext generates the YAML context for validation errors
func GenerateYamlContext(err gojsonschema.ResultError, updatedDescription string) string {
	path := err.Field()
	yaml := YamlContextBuilder{}

	parts := strings.Split(path, ".")
	if len(parts) > 0 {
		for i, part := range parts {
			yaml.WriteYamlKey(part, i)

			if i == len(parts)-1 {
				// Add offending property key
				if err.Type() == "required" || err.Type() == "additional_property_not_allowed" {
					yaml.WriteYamlKey(err.Details()["property"], i+1)
				}

				if errMap, ok := err.Value().(map[string]interface{}); ok {
					if len(errMap) > 0 {
						yaml.WriteError(updatedDescription)
						writeMapValue(&yaml, errMap, i)
					} else {
						yaml.WriteError(updatedDescription)
					}
				} else {
					yaml.WriteString(fmt.Sprintf(" %v", err.Value()))
					yaml.WriteError(updatedDescription)
				}
			}
		}
	}

	return yaml.String()
}

func writeMapValue(contextBuilder *YamlContextBuilder, errMap map[string]interface{}, indent int) {
	// Only print if the value is not a map (i.e. a primitive type)
	if len(errMap) > 0 {
		for key, value := range errMap {
			contextBuilder.WriteYamlKey(key, indent+1)

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
