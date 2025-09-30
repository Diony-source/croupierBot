package definitions

import (
	"github.com/Diony-source/CroupierBot/internal/command"
	"github.com/bwmarrin/discordgo"
)

// --- Selam Command ---
type SelamCommand struct{}

func (c *SelamCommand) Name() string { return "selam" }
func (c *SelamCommand) Execute(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	s.ChannelMessageSend(m.ChannelID, "Selam o/ Gazino'ya hoş geldin! Bu hafta The Great Vault kasası sana çalışsın. Bol keyler!")
}

// --- Hi Command ---
type HiCommand struct{}

func (h *HiCommand) Name() string { return "hi" }
func (h *HiCommand) Execute(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	s.ChannelMessageSend(m.ChannelID, "Hey there o/ Welcome to the Gazino! May the Great Vault pay out the jackpot for you this week. Happy key running!")
}

// RegisterGreetingCommands registers the greeting commands to the command handler.
func RegisterGreetingCommands(h *command.Handler) {
	h.RegisterCommand(&SelamCommand{})
	h.RegisterCommand(&HiCommand{})
}
