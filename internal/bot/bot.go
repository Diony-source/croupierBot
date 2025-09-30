// internal/bot/bot.go
package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/Diony-source/CroupierBot/internal/config"
)

// Bot struct holds the core Discord session.
type Bot struct {
	Session *discordgo.Session
}

// New creates a new Bot instance and its Discord session.
func New() (*Bot, error) {
	s, err := discordgo.New("Bot " + config.Cfg.BotToken)
	if err != nil {
		return nil, fmt.Errorf("error creating Discord session: %w", err)
	}
	return &Bot{Session: s}, nil
}

// Start opens a websocket connection to Discord and registers handlers.
func (b *Bot) Start() error {
	fmt.Println("Registering event handlers...")
	b.registerHandlers() // YENİ EKLENEN SATIR

	// We need to specify which events we want to receive from Discord.
	// For now, we only need to know about messages in guilds.
	b.Session.Identify.Intents = discordgo.IntentsGuildMessages // YENİ EKLENEN SATIR

	fmt.Println("Attempting to connect to Discord...")
	return b.Session.Open()
}

// Stop cleanly closes the Discord session.
func (b *Bot) Stop() {
	fmt.Println("\nClosing Discord session...")
	b.Session.Close()
}