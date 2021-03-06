package ball

import (
    "math/rand"
    "github.com/bwmarrin/discordgo"
)

const COMMAND string = "8ball"
const DESCRIPTION string = "ask a question and it will be answered with a yes or no"
const COLOR int = 0xff93ac

var _ANSWERS = [...]string{
	// yes like answers
	"Yes",
	"It is certain",
	"It is decidedly so",
	"Without a doubt",
	"Yes definitely",
	"You may rely on it",
	"As I see it yes",
	"Most likely",
	"Outlook good",
	"Signs point to yes",

	// no like answers
	"No",
	"My reply is no",
	"My sources say no",
	"Don't count on it",
	"Outlook not so good",
	"As I see it no",
	"Signs point to no",
	"Not likely",
	"Very doubtful",
}

type EightBall struct {}

/* for BotCommand interface */
func (self *EightBall) Prefix() string {
    return COMMAND
}
/* for BotCommand interface */
func (self *EightBall) Description() string {
    return DESCRIPTION
}
/* for BotCommand interface */
func (self *EightBall) HelpMessage(s *discordgo.Session, m *discordgo.MessageCreate) error {
    return self.ProcessCommand(s, m, "help")
}
/* for BotCommand interface */
func (self *EightBall) ProcessCommand(s *discordgo.Session, m *discordgo.MessageCreate, args string) error {
    if len(args) == 0 {
        s.ChannelMessageSend(m.ChannelID, "No questions? Back to farming...")
    } else {
        msg_index := rand.Int() % len(_ANSWERS)
        embed := &discordgo.MessageEmbed {
            Title: "Question",
            Color: COLOR,
            Description: args,
            Fields: []*discordgo.MessageEmbedField{
                { Name: "Answer", Value: _ANSWERS[msg_index] },
            },
        }
        s.ChannelMessageSendEmbed(m.ChannelID, embed)
    }
    return nil
}
