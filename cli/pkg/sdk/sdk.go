package sdk

import (
	"errors"
	"fmt"

	"github.com/nitrictech/nitric/cli/pkg/schema"
	"github.com/spf13/afero"
)

// Enum to represent SDK languages
type Language int

const (
	Go Language = iota
	Python
	Javascript
	Typescript
)

// Convert string to Language enum
func StringToLanguage(lang string) (Language, error) {
	switch lang {
	case "go":
		return Go, nil
	case "python":
		return Python, nil
	case "javascript":
		return Javascript, nil
	case "typescript":
		return Typescript, nil
	default:
		return -1, fmt.Errorf("unsupported language: %s", lang)
	}
}

// Convert Language enum to string
func LanguageToString(lang Language) (string, error) {
	switch lang {
	case Go:
		return "go", nil
	case Python:
		return "python", nil
	case Javascript:
		return "javascript", nil
	case Typescript:
		return "typescript", nil
	default:
		return "", fmt.Errorf("unsupported language: %d", lang)
	}
}

// GenerateSDK generates SDK for the specified language
func GenerateSDKs(fs afero.Fs, appSpec schema.Application, outPath string, langs []Language) error {
	// Remove duplicates
	langSet := make(map[Language]struct{})
	for _, lang := range langs {
		langSet[lang] = struct{}{}
	}

	for lang := range langSet {
		switch lang {
		case Go:
			return GenerateGoSDK(fs, appSpec, outPath)
		case Python:
			return errors.New("python SDK generation not implemented")
		case Javascript:
			return errors.New("javascript SDK generation not implemented")
		case Typescript:
			return errors.New("typescript SDK generation not implemented")
		default:
			return fmt.Errorf("unsupported language: %d", lang)
		}
	}

	return fmt.Errorf("no valid languages provided")
}
