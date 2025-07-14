package config

import (
	"reflect"
)

type ConfigKey struct {
	Path        string
	Description string
}

func (c *Config) AllKeysWithDescriptions() []ConfigKey {
	return listKeysWithDesc(c, "", "")
}

func listKeysWithDesc(data interface{}, prefix string, parentPath string) []ConfigKey {
	var keys []ConfigKey

	rv := reflect.ValueOf(data)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Struct:
		rt := rv.Type()
		for i := 0; i < rv.NumField(); i++ {
			field := rt.Field(i)
			if field.PkgPath != "" || field.Tag.Get("mapstructure") == "-" {
				continue
			}

			fieldName := field.Tag.Get("mapstructure")
			if fieldName == "" {
				fieldName = field.Name
			}

			fullKey := fieldName
			if prefix != "" {
				fullKey = prefix + "." + fieldName
			}

			fieldValue := rv.Field(i)
			subKeys := listKeysWithDesc(fieldValue.Interface(), fullKey, fullKey)

			if len(subKeys) == 0 {
				// Leaf node
				keys = append(keys, ConfigKey{
					Path:        fullKey,
					Description: field.Tag.Get("desc"),
				})
			} else {
				keys = append(keys, subKeys...)
			}
		}
	case reflect.Map:
		// Handle map types if needed
	}

	return keys
}
