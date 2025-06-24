package schema

import (
	"fmt"
	"strings"

	"github.com/spf13/afero"
	"github.com/xeipuuv/gojsonschema"
)

func LoadFromFile(fs afero.Fs, path string) (*Application, error) {
	if exists, err := afero.Exists(fs, path); err != nil {
		return nil, fmt.Errorf("nitric application file could not be loaded at path: %s", path)
	} else if !exists {
		return nil, fmt.Errorf("nitric application file not found at path: %s", path)
	}

	fileData, err := afero.ReadFile(fs, path)
	if err != nil {
		return nil, fmt.Errorf("Error opening nitric application file at path %s: %s", path, err)
	}

	var appSpec *Application
	var results *gojsonschema.Result

	if strings.HasSuffix(path, ".json") {
		appSpec, results, err = ApplicationFromJson(string(fileData))
	} else if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
		appSpec, results, err = ApplicationFromYaml(string(fileData))
	} else {
		return nil, fmt.Errorf("nitric application file must be a .json or .yaml/.yml file: %s", path)
	}

	if err != nil {
		return nil, fmt.Errorf("error parsing nitric application file %s: %s", path, err)
	}

	if results != nil && !results.Valid() {
		errs := ""
		for _, err := range results.Errors() {
			errs += fmt.Sprintf(" - %s\n", err)
		}
		return nil, fmt.Errorf("invalid nitric application file %s:\n%s", path, errs)
	}

	// Perform additional validation checks on the application
	if err := appSpec.IsValid(); err != nil {
		return nil, fmt.Errorf("invalid nitric application file %s:\n%s", path, err)
	}

	return appSpec, nil
}
