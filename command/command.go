package command

import (
    "github.com/bwmarrin/discordgo"
)

type BotCommand interface {
    Prefix() string
    Description() string
    HelpMessage(s *discordgo.Session, m *discordgo.MessageCreate) error
    ProcessCommand(s *discordgo.Session, m *discordgo.MessageCreate, args string) error
}

type CommandError struct {
    Reason string
}
func (self *CommandError) Error() string {
    return self.Reason
}
