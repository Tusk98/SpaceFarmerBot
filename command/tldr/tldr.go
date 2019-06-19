package tldr

import (
    "bytes"
    "bufio"
    "fmt"
    "net/http"
    "strings"
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

	scanner := bufio.NewScanner(resp.Body)
    var title string
    var description bytes.Buffer
    var fields []*discordgo.MessageEmbedField
    var curr_field *discordgo.MessageEmbedField
    for scanner.Scan() {
        line := scanner.Text()
        if strings.HasPrefix(line, "#") {
            title = line[2:]
        } else if strings.HasPrefix(line, ">") {
            description.WriteString(line[2:])
            description.WriteRune('\n')
        } else if strings.HasPrefix(line, "-") {
            curr_field = &discordgo.MessageEmbedField { Name: line[2:] }
        } else if strings.HasPrefix(line, "`") {
            curr_field.Value = fmt.Sprintf("``%s``", line)
            fields = append(fields, curr_field)
        }
    }
    embed := &discordgo.MessageEmbed {
        Title: title,
        Color: COLOR,
        Description: description.String(),
        Fields: fields,
    }
    s.ChannelMessageSendEmbed(m.ChannelID, embed)
    return nil
}

