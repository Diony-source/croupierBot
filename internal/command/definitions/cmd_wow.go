package definitions

import (
	"fmt"

	"github.com/Diony-source/CroupierBot/internal/command"
	"github.com/Diony-source/CroupierBot/internal/wow"
	"github.com/bwmarrin/discordgo"
)

// --- Affix Command ---
type AffixCommand struct{}

func (c *AffixCommand) Name() string { return "affix" }

func (c *AffixCommand) Execute(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	// Let the user know we are fetching the data.
	msg, err := s.ChannelMessageSend(m.ChannelID, "Fetching current Mythic+ affixes...")
	if err != nil {
		fmt.Println("Error sending initial message:", err)
		return
	}

	affixes, err := wow.GetCurrentAffixes()
	if err != nil {
		s.ChannelMessageEdit(msg.ChannelID, msg.ID, "Sorry, I couldn't get the affixes right now.")
		fmt.Println("Error getting affixes from API:", err)
		return
	}

	// Create a nice-looking response using Discord's "Embed" feature.
	embed := &discordgo.MessageEmbed{
		Title:       "This Week's Mythic+ Affixes",
		Description: fmt.Sprintf("**%s**", affixes.Title),
		Color:       0x2ECC71, // A nice green color
		Fields:      []*discordgo.MessageEmbedField{},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Data provided by Raider.IO",
		},
	}

	for _, affix := range affixes.AffixDetails {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   affix.Name,
			Value:  affix.Description,
			Inline: false,
		})
	}

	// --- DÜZELTME BURADA ---
	// 1. We create a slice variable that holds our embed(s).
	embedsToSend := []*discordgo.MessageEmbed{embed}

	// 2. We edit the message, passing a POINTER to our slice variable.
	s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Content: new(string),
		Embeds:  &embedsToSend, // The "&" symbol gets the memory address (pointer).
		ID:      msg.ID,
		Channel: m.ChannelID,
	})
}

// RegisterWoWCommands registers the WoW-related commands to the command handler.
func RegisterWoWCommands(h *command.Handler) {
	h.RegisterCommand(&AffixCommand{})
}
