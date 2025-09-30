// internal/bot/bot.go
package bot

import (
	"fmt"

	"github.com/Diony-source/CroupierBot/internal/config"
	"github.com/bwmarrin/discordgo"
)

// Bot struct holds all the core components of our bot.
type Bot struct {
	Session *discordgo.Session
	// In the future, we will add our command handler here.
}

// New creates a new Bot instance.
func New() (*Bot, error) {
	// Create a new Discord session using the provided bot token.
	s, err := discordgo.New("Bot " + config.Cfg.BotToken)
	if err != nil {
		return nil, fmt.Errorf("error creating Discord session: %w", err)
	}

	b := &Bot{
		Session: s,
	}

	return b, nil
}

// Start opens a websocket connection to Discord and begins listening.
func (b *Bot) Start() error {
	fmt.Println("Attempting to connect to Discord...")
	err := b.Session.Open()
	if err != nil {
		return fmt.Errorf("error opening connection: %w", err)
	}
	return nil
}

// Stop cleanly closes the Discord session.
func (b *Bot) Stop() {
	fmt.Println("Closing Discord session...")
	b.Session.Close()
}
