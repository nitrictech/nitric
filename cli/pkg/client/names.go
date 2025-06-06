package sdk

import (
	"fmt"
	"regexp"
	"strings"
)

type ResourceNameNormalizer struct {
	original   string
	normalized string
}

func NewResourceNameNormalizer(name string) (ResourceNameNormalizer, error) {
	if len(name) == 0 {
		return ResourceNameNormalizer{}, fmt.Errorf("resource name cannot be empty")
	}

	if !regexp.MustCompile(`^[a-zA-Z]`).MatchString(name) {
		return ResourceNameNormalizer{}, fmt.Errorf("resource name must start with a letter")
	}

	normalized := regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(name, "_")
	normalized = strings.ToLower(normalized)

	return ResourceNameNormalizer{
		original:   name,
		normalized: normalized,
	}, nil
}

func (n *ResourceNameNormalizer) Parts() []string {
	parts := strings.Split(n.normalized, "_")

	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] == "" {
			parts = append(parts[:i], parts[i+1:]...)
		}
	}

	return parts
}

// Unmodified returns the original name as it appears in the Nitric application spec
func (n *ResourceNameNormalizer) Unmodified() string {
	return n.original
}

// PascalCase returns the name in PascalCase, aka UpperCamelCase, e.g. "MyBucket"
func (n *ResourceNameNormalizer) PascalCase() string {
	return n.toCamelCase(true)
}

// CamelCase returns the name in camelCase, e.g. "myBucket"
func (n *ResourceNameNormalizer) CamelCase() string {
	return n.toCamelCase(false)
}

// SnakeCase returns the name in snake_case, e.g. "my_bucket"
func (n *ResourceNameNormalizer) SnakeCase() string {
	return strings.Join(n.Parts(), "_")
}

// KebabCase returns the name in kebab-case, e.g. "my-bucket"
func (n *ResourceNameNormalizer) KebabCase() string {
	return strings.Join(n.Parts(), "-")
}

// capitalizeFirstLetter capitalizes the first letter of the string, e.g. "my_bucket" -> "My_bucket"
func capitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// toCamelCase converts the name to camelCase or PascalCase (aka UpperCamelCase)
func (n *ResourceNameNormalizer) toCamelCase(upperCamelCase bool) string {
	var result strings.Builder
	for i, part := range n.Parts() {
		if part == "" {
			continue
		}

		if i == 0 && !upperCamelCase {
			result.WriteString(part)
			continue
		}

		result.WriteString(capitalizeFirstLetter(part))
	}
	return result.String()
}
