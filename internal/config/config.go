package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config struct now includes fields for Blizzard API credentials.
type Config struct {
	BotToken             string `json:"BotToken"`
	CommandPrefix        string `json:"CommandPrefix"`
	BlizzardClientID     string `json:"BlizzardClientID"`     // YENİ
	BlizzardClientSecret string `json:"BlizzardClientSecret"` // YENİ
}

// Cfg is a global variable that will hold our settings once loaded.
var Cfg Config

// LoadConfig will now also read and validate the new Blizzard credentials.
func LoadConfig() error {
	file, err := os.ReadFile("./config.json")
	if err != nil {
		return fmt.Errorf("could not read config file: %w", err)
	}

	err = json.Unmarshal(file, &Cfg)
	if err != nil {
		return fmt.Errorf("could not parse config file: %w", err)
	}

	// Validate that all necessary tokens and keys are present.
	if Cfg.BotToken == "" {
		return fmt.Errorf("botToken is missing or empty in config file")
	}
	if Cfg.BlizzardClientID == "" {
		return fmt.Errorf("blizzardClientID is missing or empty in config file")
	}
	if Cfg.BlizzardClientSecret == "" {
		return fmt.Errorf("blizzardClientSecret is missing or empty in config file")
	}

	return nil
}
