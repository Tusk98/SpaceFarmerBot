package foundation

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

const COMMAND string = "foundation"
const DESCRIPTION string = "SCP Foundation section of the Bot"
const WIKI string = "http://www.scp-wiki.net/scp-"
const COLOR int = 0xff93ac

//var link string = ""

type FoundationCommand struct{}

func randomInt() int {
	rand.Seed(time.Now().UnixNano())
	return 2 + rand.Intn(5000)
}

func (self *FoundationCommand) Prefix() string {
	return COMMAND
}

func (self *FoundationCommand) Description() string {
	return DESCRIPTION
}

func (self *FoundationCommand) HelpMessage(s *discordgo.Session, m *discordgo.MessageCreate) error {
	return self.ProcessCommand(s, m, "help")
}

func (self *FoundationCommand) ProcessCommand(s *discordgo.Session, m *discordgo.MessageCreate, args string) error {
	fmt.Println("INSIDE PROCESS")
	if len(args) == 0 {
		s.ChannelMessageSend(m.ChannelID, "No SCP Requests? Back to farming...")
		fmt.Println("NOT DAILY")
	} else {
		number := randomInt()
		scp := "SCP-" + strconv.Itoa(number)
		fullscp := WIKI + strconv.Itoa(number)

		embed := &discordgo.MessageEmbed{
			Title:       "Request number",
			Color:       COLOR,
			Description: scp,
			Fields: []*discordgo.MessageEmbedField{
				{Name: "SCP File", Value: fullscp},
			},
		}
		fmt.Println("ELSE")
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
	}
	return nil
}
