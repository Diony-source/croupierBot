package definitions

import (
	"fmt"
	"strings"
	"time" // NEW: We need this package for the timestamp.

	"github.com/Diony-source/CroupierBot/internal/command"
	"github.com/Diony-source/CroupierBot/internal/wow"
	"github.com/bwmarrin/discordgo"
)

// --- Affix Command ---
type AffixCommand struct{}

func (c *AffixCommand) Name() string { return "affix" }

func (c *AffixCommand) Execute(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
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

	// --- EMBED IMPROVEMENTS START HERE ---

	// Determine the color and a more descriptive title based on the main affix.
	var embedColor int
	var weekType string
	if strings.Contains(affixes.Title, "Tyrannical") {
		embedColor = 0xC41E3A // A deep red color for Tyrannical
		weekType = "Tyrannical Week"
	} else {
		embedColor = 0x3498DB // A nice blue color for Fortified
		weekType = "Fortified Week"
	}

	// Create a new, more beautiful embed.
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "CroupierBot",
			IconURL: "https://wow.zamimg.com/images/wow/icons/large/inv_relics_key_01.jpg", // M+ Keystone Icon in author line
		},
		Title:       weekType,
		Description: "Here are the active affixes for this week:",
		Color:       embedColor,
		Fields:      []*discordgo.MessageEmbedField{},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://wow.zamimg.com/images/wow/icons/large/inv_relics_key_01.jpg", // M+ Keystone Icon as thumbnail
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Data provided by Raider.IO",
		},
		Timestamp: time.Now().Format(time.RFC3339), // Adds the current time to the footer.
	}

	for _, affix := range affixes.AffixDetails {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   affix.Name,
			Value:  affix.Description,
			Inline: false,
		})
	}

	// --- EMBED IMPROVEMENTS END HERE ---

	embedsToSend := []*discordgo.MessageEmbed{embed}

	s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Content: new(string),
		Embeds:  &embedsToSend,
		ID:      msg.ID,
		Channel: m.ChannelID,
	})
}

// RegisterWoWCommands registers the WoW-related commands to the command handler.
func RegisterWoWCommands(h *command.Handler) {
	h.RegisterCommand(&AffixCommand{})
}
