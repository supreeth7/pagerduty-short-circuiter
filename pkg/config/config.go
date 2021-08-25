package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Configuration struct to store user configuration
type Config struct {
	ApiKey string `json:"api_key,omitempty"`
}

// Find returns the location of the pdcli configuration file.
func Find() (string, error) {
	//Fetches the user home directory
	homedir, err := os.UserHomeDir()

	if err != nil {
		return "", fmt.Errorf("cannot locate user home directory")
	}

	configPath := filepath.Join(homedir, ".config/pagerduty-cli/config.json")

	//Check if the path exists
	_, err = os.Stat(configPath)

	//If the config path doesn't exist
	if os.IsNotExist(err) {
		configDir, err := os.UserConfigDir()

		if err != nil {
			return configDir, err
		}

		configPath = filepath.Join(configDir, "pagerduty-cli/config.json")

		// Creates a new directory inside the .config dir
		dir := filepath.Dir(configPath)
		err = os.MkdirAll(dir, os.FileMode(0755))

		if err != nil {
			return "", err
		}
	}

	return configPath, nil
}

func Fetch() (config *Config, err error) {
	//Locate the config file
	configFile, err := Find()

	if err != nil {
		return
	}

	_, err = os.Stat(configFile)

	if os.IsNotExist(err) {
		config = &Config{}
		err = nil
		return
	}

	configData, err := ioutil.ReadFile(configFile)

	if err != nil {
		err = fmt.Errorf("can't read config file")
		return
	}

	config = &Config{}

	if len(configData) == 0 {
		err = fmt.Errorf("configuration file is empty")
		return
	}

	err = json.Unmarshal(configData, config)

	if err != nil {
		err = fmt.Errorf("error parsing config file")
		return
	}

	return
}

// Save saves the given configuration to the configuration file.
func Save(cfg *Config) error {
	file, err := Find()

	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")

	if err != nil {
		return fmt.Errorf("can't marshal config: %v", err)
	}

	err = ioutil.WriteFile(file, data, 0600)

	if err != nil {
		return fmt.Errorf("can't write file '%s': %v", file, err)
	}

	return nil
}
