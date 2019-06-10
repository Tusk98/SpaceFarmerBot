package ball

import (
    "math/rand"
    "github.com/bwmarrin/discordgo"
)

const Command string = "8ball"

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

func ProcessCommand(s *discordgo.Session, m *discordgo.MessageCreate, args string) error {
    if args == "" {
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
