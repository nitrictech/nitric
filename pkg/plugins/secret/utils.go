package secret

import (
	"fmt"
	"regexp"

	"github.com/nitric-dev/membrane/pkg/plugins"
)

const rx = "^\\w+(-\\w+)*$"

// ValidateSecretName - Validates a secret name
func ValidateSecretName(secName string) error {
	if len(secName) == 0 {
		return plugins.NewInvalidArgError("Secret name must not be blank")
	}

	match, _ := regexp.MatchString(rx, secName)

	if !match {
		return plugins.NewInvalidArgError(fmt.Sprintf("Secret name must match pattern: %s", rx))
	}

	return nil
}