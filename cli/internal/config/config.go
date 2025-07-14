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
	v   *viper.Viper `mapstructure:"-"`
	Url string       `mapstructure:"url" desc:"The base URL of the Nitric server (e.g., https://app.nitric.io)"`
}

func (c *Config) FileUsed() string {
	return c.v.ConfigFileUsed()
}

func (c *Config) Dump() string {
	custom := dump(c.v.AllSettings())

	var configMap map[string]interface{}
	err := mapstructure.Decode(c, &configMap)
	if err != nil {
		return ""
	}

	all := dump(configMap)

	b := strings.Builder{}

	for key, value := range all {
		if customValue, ok := custom[key]; ok {
			b.WriteString(fmt.Sprintf("%s: %s\n", key, customValue))
		} else {
			b.WriteString(fmt.Sprintf("%s: %s (default)\n", key, value))
		}
	}

	return b.String()
}

type ConfigValue struct {
	Key   string
	Value string
}

func dump(data map[string]any) map[string]string {
	allValues := map[string]string{}
	for key, value := range data {
		if _, ok := value.(map[string]any); ok {
			subValues := dump(value.(map[string]any))
			for k, v := range subValues {
				allValues[key+"."+k] = v
			}
		} else {
			allValues[key] = fmt.Sprintf("%v", value)
		}
	}
	return allValues
}

func (c *Config) SetValue(key, value string) error {
	return mapstructure.Decode(map[string]interface{}{key: value}, c)
}

func (c *Config) GetNitricServerUrl() *url.URL {
	nitricUrl, err := url.Parse(c.Url)
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

	c.Url = nitricUrl.String()
	return nil
}

func (c *Config) Save(global bool) error {
	var configMap map[string]interface{}
	err := mapstructure.Decode(c, &configMap)
	if err != nil {
		return fmt.Errorf("failed to decode config struct: %w", err)
	}

	for key, value := range configMap {
		c.v.Set(key, value)
	}

	var file string
	if global {
		homeConfigPath, err := HomeConfigPath()
		if err != nil {
			return err
		}
		file = filepath.Join(homeConfigPath, "config.yaml")
	} else {
		file = filepath.Join(LocalConfigPath(), "config.yaml")
	}

	// Make sure the config file exists
	if _, err := os.Stat(file); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(file), 0755)
		if err != nil {
			return fmt.Errorf("failed to create config file: %w", err)
		}
	}

	c.v.SetConfigFile(file)

	err = c.v.WriteConfig()
	if err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

func HomeConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".nitric"), nil
}

func LocalConfigPath() string {
	return ".nitric"
}

func Load(cmd *cobra.Command) (*Config, error) {
	v := viper.New()

	v.SetConfigType("yaml")

	// Search the current .nitric directory first
	v.AddConfigPath(LocalConfigPath())

	homeConfigPath, err := HomeConfigPath()
	if err != nil {
		return nil, err
	}
	v.AddConfigPath(homeConfigPath)

	c := Config{
		v: v,
		// Set any default values here
		Url: "https://app.nitric.io/",
	}

	err = v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return &c, nil
		}

		return nil, err
	}

	err = v.Unmarshal(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
