package definitions

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Diony-source/CroupierBot/internal/command"
	"github.com/Diony-source/CroupierBot/internal/wow"
	"github.com/bwmarrin/discordgo"
)

// --- Affix Command (No changes) ---
type AffixCommand struct{}

func (c *AffixCommand) Name() string { return "affix" }
func (c *AffixCommand) Execute(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	// ... (This command's code remains unchanged)
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
		embedColor = 0xC41E3A
		weekType = "Tyrannical Week"
	} else {
		embedColor = 0x3498DB
		weekType = "Fortified Week"
	}
	embed := &discordgo.MessageEmbed{Author: &discordgo.MessageEmbedAuthor{Name: "CroupierBot", IconURL: "https://wow.zamimg.com/images/wow/icons/large/inv_relics_key_01.jpg"}, Title: weekType, Description: "Here are the active affixes for this week:", Color: embedColor, Fields: []*discordgo.MessageEmbedField{}, Thumbnail: &discordgo.MessageEmbedThumbnail{URL: "https://wow.zamimg.com/images/wow/icons/large/inv_relics_key_01.jpg"}, Footer: &discordgo.MessageEmbedFooter{Text: "Data provided by Raider.IO"}, Timestamp: time.Now().Format(time.RFC3339)}
	for _, affix := range affixes.AffixDetails {
		emoji := getAffixEmoji(affix.Name)
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{Name: fmt.Sprintf("%s %s", emoji, affix.Name), Value: affix.Description, Inline: false})
	}
	embedsToSend := []*discordgo.MessageEmbed{embed}
	s.ChannelMessageEditComplex(&discordgo.MessageEdit{Content: new(string), Embeds: &embedsToSend, ID: msg.ID, Channel: m.ChannelID})
}

// --- The Ultimate Char Command (CORRECTED AND ENHANCED) ---
type CharCommand struct{}

func (c *CharCommand) Name() string { return "char" }
func (c *CharCommand) Execute(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Please use the format: `!char <character-name> <server-name>`. Example: `!char Freythemercy Silvermoon`")
		return
	}
	characterName := args[0]
	serverNameBlizz := strings.Join(args[1:], " ")
	serverNameRio := strings.Join(args[1:], "-")
	region := "eu"

	msg, _ := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Fetching a full report for %s-%s...", characterName, serverNameBlizz))

	var wg sync.WaitGroup
	var rioProfile *wow.CharacterProfileResponse
	var bnetSummary *wow.CharacterSummary
	var bnetMedia *wow.CharacterMediaSummary
	var bnetEquipment *wow.CharacterEquipment
	var rioErr, summaryErr, mediaErr, equipmentErr error

	wg.Add(4)
	go func() {
		defer wg.Done()
		rioProfile, rioErr = wow.GetCharacterProfile(characterName, serverNameRio, region)
	}()
	go func() {
		defer wg.Done()
		bnetSummary, summaryErr = wow.GetCharacterSummary(characterName, serverNameBlizz, region)
	}()
	go func() {
		defer wg.Done()
		bnetMedia, mediaErr = wow.GetCharacterMedia(characterName, serverNameBlizz, region)
	}()
	go func() {
		defer wg.Done()
		bnetEquipment, equipmentErr = wow.GetCharacterEquipment(characterName, serverNameBlizz, region)
	}()
	wg.Wait()

	if summaryErr != nil {
		s.ChannelMessageEdit(msg.ChannelID, msg.ID, "Could not fetch character data from Blizzard Armory. Please check spelling.")
		fmt.Println("Blizzard Summary Error:", summaryErr)
		return
	}

	classColor := getClassColor(bnetSummary.CharacterClass.Name)

	characterImageURL := ""
	if mediaErr == nil && bnetMedia != nil {
		for _, asset := range bnetMedia.Assets {
			if asset.Key == "main-raw" {
				characterImageURL = asset.Value
				break
			}
		}
	}

	mythicPlusScore := "N/A"
	raidProgress := "N/A"
	if rioErr == nil && rioProfile != nil {
		if len(rioProfile.MythicPlusScoresBySeason) > 0 {
			mythicPlusScore = fmt.Sprintf("%.1f", rioProfile.MythicPlusScoresBySeason[0].Scores.All)
		}
		if progress, ok := rioProfile.RaidProgression["manaforge-omega"]; ok {
			raidProgress = progress.Summary
		}
	}

	// UPDATED: Equipment parsing logic is now corrected and enhanced.
	equipmentString := ""
	if equipmentErr == nil && bnetEquipment != nil {
		var sb strings.Builder
		for _, item := range bnetEquipment.EquippedItems {
			enchant := ""
			if len(item.Enchantments) > 0 {
				// Parse the enchant name from the display string.
				enchantName := strings.Split(item.Enchantments[0].DisplayString, "|")[0]
				enchant = fmt.Sprintf(" ✨ *%s*", strings.TrimSpace(enchantName))
			}
			// Use item.Level.Value (correct) instead of item.Item.Level.Value (incorrect).
			sb.WriteString(fmt.Sprintf("**%s**: `%.0f`%s\n", item.Slot.Name, item.Level.Value, enchant))
		}
		equipmentString = sb.String()
	}
	if equipmentString == "" {
		equipmentString = "Could not fetch equipment details."
	}

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name: fmt.Sprintf("%s - %s (%s)", bnetSummary.Name, bnetSummary.Realm.Slug, strings.ToUpper(region)),
		},
		Color: classColor,
		Image: &discordgo.MessageEmbedImage{
			URL: characterImageURL,
		},
		Fields: []*discordgo.MessageEmbedField{
			{Name: "iLvl (Equipped/Avg)", Value: fmt.Sprintf("%.0f / %.0f", bnetSummary.EquippedItemLevel, bnetSummary.AverageItemLevel), Inline: true},
			{Name: "Class & Spec", Value: fmt.Sprintf("%s (%s)", bnetSummary.CharacterClass.Name, bnetSummary.ActiveSpec.Name), Inline: true},
			{Name: "M+ Score", Value: mythicPlusScore, Inline: false},
			{Name: "Raid Progress", Value: raidProgress, Inline: false},
			{Name: "Equipment", Value: equipmentString, Inline: false},
		},
		Footer:    &discordgo.MessageEmbedFooter{Text: "Data from Blizzard & Raider.IO"},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	embedsToSend := []*discordgo.MessageEmbed{embed}
	s.ChannelMessageEditComplex(&discordgo.MessageEdit{Content: new(string), Embeds: &embedsToSend, ID: msg.ID, Channel: m.ChannelID})
}

// --- Helper Functions (No changes) ---
func getClassColor(className string) int {
	// ... (This function remains unchanged)
	switch className {
	case "Death Knight":
		return 0xC41F3B
	case "Demon Hunter":
		return 0xA330C9
	case "Druid":
		return 0xFF7D0A
	case "Hunter":
		return 0xABD473
	case "Mage":
		return 0x69CCF0
	case "Monk":
		return 0x00FF96
	case "Paladin":
		return 0xF58CBA
	case "Priest":
		return 0xFFFFFF
	case "Rogue":
		return 0xFFF569
	case "Shaman":
		return 0x0070DE
	case "Warlock":
		return 0x9482C9
	case "Warrior":
		return 0xC79C6E
	case "Evoker":
		return 0x33937F
	default:
		return 0xAAAAAA
	}
}
func getAffixEmoji(affixName string) string {
	// ... (This function remains unchanged)
	switch affixName {
	case "Fortified":
		return "🛡️"
	case "Tyrannical":
		return "👑"
	case "Spiteful":
		return "👻"
	case "Raging":
		return "😡"
	case "Bolstering":
		return "💪"
	case "Sanguine":
		return "🩸"
	case "Bursting":
		return "💥"
	case "Volcanic":
		return "🌋"
	case "Storming":
		return "🌪️"
	case "Explosive":
		return "💣"
	default:
		return "🔹"
	}
}

// --- Command Registration (No changes) ---
func RegisterWoWCommands(h *command.Handler) {
	h.RegisterCommand(&AffixCommand{})
	h.RegisterCommand(&CharCommand{})
}
