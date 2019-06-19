package tldr

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "github.com/bwmarrin/discordgo"
)

const COMMAND string = "tldr"
const DESCRIPTION string = "fetches practical examples of linux commands"
const COLOR int = 0xff93ac

const TLDR_INDEX = "https://raw.githubusercontent.com/tldr-pages/tldr/master/pages"

type TldrCommand struct {}

/* for BotCommand interface */
func (self *TldrCommand) Prefix() string {
    return COMMAND
}
/* for BotCommand interface */
func (self *TldrCommand) Description() string {
    return DESCRIPTION
}
/* for BotCommand interface */
func (self *TldrCommand) HelpMessage(s *discordgo.Session, m *discordgo.MessageCreate) error {
    embed := &discordgo.MessageEmbed {
        Title: fmt.Sprintf("%s usage", self.Prefix()),
        Color: COLOR,
        Description: fmt.Sprintf("Usage: %s command\n%s", self.Prefix(), self.Description()),
    }
    s.ChannelMessageSendEmbed(m.ChannelID, embed)
    return nil
}
/* for BotCommand interface */
func (self *TldrCommand) ProcessCommand(s *discordgo.Session, m *discordgo.MessageCreate, args string) error {
    /* first check if its a linux specific command */
    url := fmt.Sprintf("%s/linux/%s.md", TLDR_INDEX, args)
    resp, err := http.Get(url)

    /* check if it is a common command, if not linux-specific */
    if resp.StatusCode != 200 {
        url := fmt.Sprintf("%s/common/%s.md", TLDR_INDEX, args)
        resp, err = http.Get(url)
    }

    if err != nil {
        return err
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```md\n%s```", body))
    return nil
}

