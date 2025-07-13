package config

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/nitrictech/nitric/cli/internal/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	v               *viper.Viper
	NitricServerUrl string `mapstructure:"url"`
}

func (c *Config) FileUsed() string {
	return c.v.ConfigFileUsed()
}

func (c *Config) Dump() string {
	all := c.v.AllSettings()

	allLines := []string{}

	for key, value := range all {
		allLines = append(allLines, fmt.Sprintf("%s: %v", key, value))
	}

	return strings.Join(allLines, "\n")
}

func (c *Config) SetValue(key, value string) error {
	return mapstructure.Decode(map[string]interface{}{key: value}, c)
}

func (c *Config) GetNitricServerUrl() *url.URL {
	nitricUrl, err := url.Parse(c.NitricServerUrl)
	if err != nil {
		fmt.Printf("Error parsing %s server url from config, using default: %v\n", version.ProductName, err)
		return &url.URL{
			Scheme: "https",
			Host:   "app.nitric.io",
		}
	}

	return nitricUrl
}

func (c *Config) SetNitricServerUrl(newUrl string) error {
	nitricUrl, err := url.Parse(newUrl)
	if err != nil {
		return err
	}

	c.NitricServerUrl = nitricUrl.String()
	return nil
}

func (c *Config) Save() error {
	var configMap map[string]interface{}
	err := mapstructure.Decode(c, &configMap)
	if err != nil {
		return fmt.Errorf("failed to decode config struct: %w", err)
	}

	for key, value := range configMap {
		c.v.Set(key, value)
	}

	err = c.v.WriteConfig()
	if err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

func Load(cmd *cobra.Command) (*Config, error) {
	v := viper.New()
	v.SetDefault("url", "https://app.nitric.io/")

	v.SetConfigType("yaml")

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	} else {
		v.AddConfigPath(filepath.Join(home, ".nitric"))
	}

	// Search the current .nitric directory first
	v.AddConfigPath(".nitric")

	err = v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	c := &Config{v: v}

	err = v.Unmarshal(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
