package foundation

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Tusk98/SpaceFarmerBot/command"
	"github.com/bwmarrin/discordgo"
	"github.com/pelletier/go-toml"
)

const COMMAND string = "foundation"
const DESCRIPTION string = "SCP Foundation section of the Bot"
const WIKI string = "http://www.scp-wiki.net/scp-"
const COLOR int = 0xff93ac

var link string = ""

func randomInt() int {
	rand.Seed(time.Now().UnixNano())
	return 2 + rand.Intn(5000)
}

type SCP struct{}

func (self *SCP) Prefix() string {
	return COMMAND
}

func (self *SCP) Description() string {
	return DESCRIPTION
}

func (self *SCP) HelpMessage(s *discordgo.Session, m *discordgo.MessageCreate) error {
	return self.ProcessCommand(s, m, "help")
}

func (self *SCP) ProcessCommand(s *discordgo.Session, m *discordgo.MessageCreate, args string) error {
	fmt.Println("INSIDE PROCESS")
	if (strings.TrimSpace(args)) != "daily" {
		s.ChannelMessageSend(m.ChannelID, "No SCP Requests? Back to farming...")
		fmt.Println("NOT DAILY")
	} else {
		number := randomInt()
		scp := "SCP-" + strconv.Itoa(number)
		fullscp := WIKI + strconv.Itoa(number)

		embed := &discordgo.MessageEmbed{
			Title:       scp,
			Color:       COLOR,
			Description: args,
			Fields: []*discordgo.MessageEmbedField{
				{Name: "SCP File", Value: link},
			},
		}
		fmt.Println("ELSE")
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
	}
	return nil
}
