package bot

import (
	"github.com/Diony-source/CroupierBot/internal/command/definitions"
	"github.com/bwmarrin/discordgo"
)

// registerCommands registers all command definitions to the command handler.
func (b *Bot) registerCommands() {
	definitions.RegisterGreetingCommands(b.CmdHandler)
	// Add other command groups here in the future
}

// registerHandlers registers all the event handlers for the bot.
func (b *Bot) registerHandlers() {
	b.Session.AddHandler(b.onMessageCreate)
	// Add more handlers here in the future
}

// onMessageCreate is called when a message is created.
// It now routes valid commands to the command handler.
func (b *Bot) onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	// Instead of printing, we now route the message.
	b.CmdHandler.Route(s, m)
}
