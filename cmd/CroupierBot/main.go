package main

import (
	"fmt"
	// This import path is based on your go.mod file.
	"github.com/Diony-source/CroupierBot/internal/config"
)

func main() {
	// Start by loading the configuration from the config.json file.
	err := config.LoadConfig()
	if err != nil {
		// If there is an error, print it and exit the program.
		// This is crucial to prevent the bot from running with a bad configuration.
		fmt.Println("Error loading configuration:", err)
		return
	}

	// If loading is successful, print a confirmation message.
	fmt.Println("Configuration loaded successfully. Bot is starting...")

	// --- Bot startup logic will go here in the future ---
}