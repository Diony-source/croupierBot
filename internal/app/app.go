// internal/app/app.go
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Diony-source/CroupierBot/internal/bot"
	"github.com/Diony-source/CroupierBot/internal/config"
)

// Run initializes and runs the application.
func Run() error {
	// 1. Load configuration.
	if err := config.LoadConfig(); err != nil {
		return fmt.Errorf("error loading configuration: %w", err)
	}

	// 2. Create a new bot instance.
	bot, err := bot.New()
	if err != nil {
		return fmt.Errorf("error creating the bot: %w", err)
	}

	// 3. Start the bot's session.
	if err := bot.Start(); err != nil {
		return fmt.Errorf("error starting the bot: %w", err)
	}

	// 4. Wait for a shutdown signal.
	fmt.Println("CroupierBot is now running. Press CTRL+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// 5. Cleanly shut down the bot's session.
	bot.Stop()

	return nil
}
