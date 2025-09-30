package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config struct holds the configuration loaded from config.json.
// The json tags are used to map the JSON keys to the struct fields.
type Config struct {
	BotToken      string `json:"BotToken"`
	CommandPrefix string `json:"CommandPrefix"`
}

// Cfg is a global variable that will hold our settings once loaded.
// Other packages can access the configuration through this variable (e.g., config.Cfg.BotToken).
var Cfg Config

// LoadConfig finds and reads the config.json file from the project's root,
// then unmarshals its content into the Cfg variable.
func LoadConfig() error {
	// Read the config.json file from the root directory.
	file, err := os.ReadFile("./config.json")
	if err != nil {
		return fmt.Errorf("could not read config file: %w", err)
	}

	// Unmarshal the JSON data into the Config struct.
	err = json.Unmarshal(file, &Cfg)
	if err != nil {
		return fmt.Errorf("could not parse config file: %w", err)
	}

	// Validate that the token is not empty, return an error if it is.
	if Cfg.BotToken == "" {
		return fmt.Errorf("botToken is missing or empty in config file")
	}

	// If everything is fine, return nil (no error).
	return nil
}
