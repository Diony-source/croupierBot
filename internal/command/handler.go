package command

import (
	"strings"

	"github.com/Diony-source/CroupierBot/internal/config"
	"github.com/bwmarrin/discordgo"
)

type ICommand interface {
	Name() string
	Execute(s *discordgo.Session, m *discordgo.MessageCreate, args []string)
}

type Handler struct {
	Commands map[string]ICommand
}

func NewHandler() *Handler {
	return &Handler{
		Commands: make(map[string]ICommand),
	}
}

func (h *Handler) RegisterCommand(cmd ICommand) {
	h.Commands[cmd.Name()] = cmd
}

func (h *Handler) Route(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !strings.HasPrefix(m.Content, config.Cfg.CommandPrefix) {
		return
	}
	parts := strings.Fields(strings.TrimPrefix(m.Content, config.Cfg.CommandPrefix))
	commandName := strings.ToLower(parts[0])
	args := parts[1:]
	if cmd, exists := h.Commands[commandName]; exists {
		cmd.Execute(s, m, args)
	}
}
