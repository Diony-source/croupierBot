// internal/bot/handlers.go
package bot

import "github.com/bwmarrin/discordgo"

// registerHandlers registers all the event handlers for the bot.
func (b *Bot) registerHandlers() {
	// Add the onMessageCreate function as a handler for the MessageCreate event.
	b.Session.AddHandler(b.onMessageCreate)

	// Add more handlers here in the future (e.g., for when a user joins).
}

// onMessageCreate is called every time a new message is created on any channel.
func (b *Bot) onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself to prevent infinite loops.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// For now, we will just print the message to the console to see if it works.
	// Later, this will be passed to our command handler.
	fmt.Printf("Message Received: GuildID=%s, ChannelID=%s, Author=%s, Content='%s'\n",
		m.GuildID, m.ChannelID, m.Author.Username, m.Content)
}