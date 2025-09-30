// cmd/CroupierBot/main.go
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Diony-source/CroupierBot/internal/bot"
	"github.com/Diony-source/CroupierBot/internal/config"
)

func main() {
	// 1. Load configuration
	if err := config.LoadConfig(); err != nil {
		fmt.Println("Error loading configuration:", err)
		return
	}

	// 2. Create a new bot instance
	krupiye, err := bot.New()
	if err != nil {
		fmt.Println("Error creating the bot:", err)
		return
	}

	// 3. Start the bot (connect to Discord)
	if err := krupiye.Start(); err != nil {
		fmt.Println("Error starting the bot:", err)
		return
	}

	fmt.Println("CroupierBot is now running. Press CTRL+C to exit.")

	// 4. Wait for a shutdown signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// 5. Cleanly shut down the bot
	krupiye.Stop()
}