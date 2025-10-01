package definitions

import (
	"fmt"
	"strings"
	"time"

	"github.com/Diony-source/CroupierBot/internal/command"
	"github.com/Diony-source/CroupierBot/internal/wow"
	"github.com/bwmarrin/discordgo"
)

// --- Affix Komutu ---
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
	var embedColor int
	var weekType string
	if strings.Contains(affixes.Title, "Tyrannical") {
		embedColor = 0xC41E3A // Tyrannical için kırmızı
		weekType = "Tyrannical Week"
	} else {
		embedColor = 0x3498DB // Fortified için mavi
		weekType = "Fortified Week"
	}
	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{Name: "CroupierBot", IconURL: "https://wow.zamimg.com/images/wow/icons/large/inv_relics_key_01.jpg"},
		Title:       weekType,
		Description: "Here are the active affixes for this week:",
		Color:       embedColor,
		Fields:      []*discordgo.MessageEmbedField{},
		Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: "https://wow.zamimg.com/images/wow/icons/large/inv_relics_key_01.jpg"},
		Footer:      &discordgo.MessageEmbedFooter{Text: "Data provided by Raider.IO"},
		Timestamp:   time.Now().Format(time.RFC3339),
	}
	for _, affix := range affixes.AffixDetails {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{Name: affix.Name, Value: affix.Description, Inline: false})
	}
	embedsToSend := []*discordgo.MessageEmbed{embed}
	s.ChannelMessageEditComplex(&discordgo.MessageEdit{Content: new(string), Embeds: &embedsToSend, ID: msg.ID, Channel: m.ChannelID})
}

// --- Rio Komutu ---
type RioCommand struct{}

func (c *RioCommand) Name() string { return "rio" }

func (c *RioCommand) Execute(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Please use the format: `!rio <character-name> <server-name>`. For example: `!rio Methodjosh Twisting-Nether`")
		return
	}
	characterName := args[0]
	serverName := args[1]

	msg, _ := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Fetching Raider.IO profile for %s-%s...", characterName, serverName))

	profile, err := wow.GetCharacterProfile(characterName, serverName, "eu") // Şimdilik bölgeyi EU varsayıyoruz
	if err != nil {
		s.ChannelMessageEdit(msg.ChannelID, msg.ID, "Could not find character. Please check the spelling of the character and server name (e.g., `Twisting-Nether`).")
		fmt.Println("Error getting character profile:", err)
		return
	}

	mythicPlusScore := "N/A"
	if profile.MythicPlusScores != nil {
		mythicPlusScore = fmt.Sprintf("%.2f", profile.MythicPlusScores.All)
	}

	raidProgress := "N/A"
	if amirdrassil, ok := profile.RaidProgression["amirdrassil-the-dreams-hope"]; ok {
		raidProgress = amirdrassil.Summary
	}

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    fmt.Sprintf("%s - %s %s", profile.Name, profile.Race, profile.Class),
			IconURL: s.State.User.AvatarURL(""),
		},
		Title: "Raider.IO Profile",
		URL:   profile.ProfileURL,
		Color: 0xff8000, // Raider.IO turuncusu
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: profile.ThumbnailURL,
		},
		Fields: []*discordgo.MessageEmbedField{
			{Name: "M+ Score", Value: mythicPlusScore, Inline: true},
			{Name: "Spec", Value: profile.ActiveSpecName, Inline: true},
			{Name: "Raid Progress (Amirdrassil)", Value: raidProgress, Inline: false},
		},
		Footer:    &discordgo.MessageEmbedFooter{Text: "Data provided by Raider.IO"},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	embedsToSend := []*discordgo.MessageEmbed{embed}
	s.ChannelMessageEditComplex(&discordgo.MessageEdit{Content: new(string), Embeds: &embedsToSend, ID: msg.ID, Channel: m.ChannelID})
}

// --- Komut Kaydı ---
func RegisterWoWCommands(h *command.Handler) {
	h.RegisterCommand(&AffixCommand{})
	h.RegisterCommand(&RioCommand{}) // Yeni komutu kaydet
}
