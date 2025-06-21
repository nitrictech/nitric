package config

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/spf13/viper"
)

const (
	ConfigFile = ".cli_config"
	ApiUrlKey  = "api.url"
	EnvPrefix  = "NITRIC"
)

var allConfigKeys = []string{ApiUrlKey}

func GetAllConfigItems() map[string]string {
	items := make(map[string]string)
	for _, key := range allConfigKeys {
		items[key] = viper.GetString(key)
	}
	return items
}

func GetValue(key string) (string, error) {
	if slices.Contains(allConfigKeys, key) {
		return viper.GetString(key), nil
	}

	return "", fmt.Errorf("invalid config option %s", key)
}

func SetValue(key string, value string) error {
	if slices.Contains(allConfigKeys, key) {
		viper.Set(key, value)
		return nil
	}

	return fmt.Errorf("invalid config option %s", key)
}

func GetApiUrl() string {
	return viper.GetString(ApiUrlKey)
}

func SetApiUrl(url string) {
	viper.Set(ApiUrlKey, url)
}

func Save() error {
	if err := viper.WriteConfig(); err != nil {
		// If config file doesn't exist, create it
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := viper.SafeWriteConfig(); err != nil {
				return fmt.Errorf("error creating config file: %v", err)
			}
		} else {
			return fmt.Errorf("error writing config: %v", err)
		}
	}

	return nil
}

func setDefaults() {
	viper.SetDefault(ApiUrlKey, "https://app.nitric.io/api/")
}

// Load loads the config from the file or the home directory
func Load(file string) error {
	if file != "" {
		// Use config file from the flag.
		viper.SetConfigFile(file)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		configDir := filepath.Join(home, ".nitric")

		// Search config in home directory with name ".nitric" (without extension).
		viper.AddConfigPath(configDir)
		// Search config in the current directory
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(ConfigFile)
	}

	setDefaults()

	viper.SetEnvPrefix(EnvPrefix)
	viper.AutomaticEnv() // read in environment variables that match

	return viper.ReadInConfig()
}
