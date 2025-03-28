package config

import "github.com/mitchellh/mapstructure"

type TagsConfig struct {
	Tags map[string]string `mapstructure:"tags"`
}

func GetTagsConfigFromAttributes(attributes map[string]interface{}) (*TagsConfig, error) {
	config := new(TagsConfig)

	err := mapstructure.Decode(attributes, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
