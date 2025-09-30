package bot

import (
	"fmt"

	"github.com/Diony-source/CroupierBot/internal/command"
	"github.com/Diony-source/CroupierBot/internal/config"
	"github.com/bwmarrin/discordgo"
)

// Bot struct now holds the command handler.
type Bot struct {
	Session    *discordgo.Session
	CmdHandler *command.Handler
}

// New creates a new Bot instance, now with a command handler.
func New() (*Bot, error) {
	s, err := discordgo.New("Bot " + config.Cfg.BotToken)
	if err != nil {
		return nil, fmt.Errorf("error creating Discord session: %w", err)
	}

	return &Bot{
		Session:    s,
		CmdHandler: command.NewHandler(),
	}, nil
}

// Start now registers commands and handlers before connecting.
func (b *Bot) Start() error {
	fmt.Println("Registering commands...")
	b.registerCommands()

	fmt.Println("Registering event handlers...")
	b.registerHandlers()

	// We must specify the Intents (what our bot wants to listen to).
	b.Session.Identify.Intents = discordgo.IntentsGuildMessages

	fmt.Println("Attempting to connect to Discord...")
	return b.Session.Open()
}

// Stop cleanly closes the Discord session.
func (b *Bot) Stop() {
	fmt.Println("\nClosing Discord session...")
	b.Session.Close()
}
