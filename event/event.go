package event

import (
    "fmt"
    "strings"
    "github.com/pelletier/go-toml"
    "github.com/bwmarrin/discordgo"
)

const COMMAND string = "event"
const DESCRIPTION string = "creates and tracks events"

const COLOR int = 0xff93ac

type Event struct {
    Name string
    Date string
    Description string
}

type EventInfoNotFoundError struct {
    reason string
}
func (self *EventInfoNotFoundError) Error() string {
    return self.reason
}

func HelpMessage(s *discordgo.Session, m *discordgo.MessageCreate) error {
    embed := &discordgo.MessageEmbed {
        Title: "event usage",
        Color: COLOR,
        Description: fmt.Sprintf("Usage: event [OPTIONS]\n%s", DESCRIPTION),
        Fields: []*discordgo.MessageEmbedField{
            { Name: "new\n```name = \"...\"\ndate = \"...\"\ndescription = \"...\"```",
                    Value: "creates a new event with the specified information" },
            { Name: "remove name", Value: "removes the given event" },
        },
    }
    s.ChannelMessageSendEmbed(m.ChannelID, embed)
    return nil
}

func ProcessCommand(s *discordgo.Session, m *discordgo.MessageCreate, args string) error {
    slice_ind := strings.IndexRune(args, '`')
    cmd := args
    xs := ""
    if slice_ind != -1 {
        cmd = strings.TrimSpace(args[:slice_ind])
        xs = args[slice_ind:]
    }
    fmt.Printf("cmd: \"%s\"\nxs: \"%s\"\n", cmd, xs)

    switch cmd {
        case "": {
            s.ChannelMessageSend(m.ChannelID, "TODO")
        }
        case "help": return HelpMessage(s, m)
        case "new": return newEvent(s, m, xs)
        case "remove": return removeEvent(s, m, xs)
        default: return HelpMessage(s, m)
    }
    return nil
}

func newEvent(s *discordgo.Session, m *discordgo.MessageCreate, args string) error {
    const DISCORD_CODE_ENCLOSURE string = "```"

    event_info := args
    event_str_start := strings.Index(event_info, DISCORD_CODE_ENCLOSURE)
    if event_str_start == -1 {
        return &EventInfoNotFoundError { "Could not find ```" }
    }
    event_info = event_info[event_str_start + len(DISCORD_CODE_ENCLOSURE):]

    event_str_end := strings.Index(event_info, DISCORD_CODE_ENCLOSURE)
    if event_str_end == -1 {
        return &EventInfoNotFoundError { "Could not find matching pair of ```" }
    }
    event_info = event_info[:event_str_end]

    var event Event
    if err := toml.Unmarshal([]byte(event_info), &event); err != nil {
        return err
    }

    embed := &discordgo.MessageEmbed {
        Title: event.Name,
        Color: COLOR,
        Description: event.Description,
        Fields: []*discordgo.MessageEmbedField{
            { Name: "Date", Value: event.Date },
            { Name: "Going", Value: "-" },
        },
    }
    msg, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
    if err != nil {
        return err
    }
    s.MessageReactionAdd(msg.ChannelID, msg.ID, "\u2705")
    s.MessageReactionAdd(msg.ChannelID, msg.ID, "\u274E")
    return nil
}

func removeEvent(s *discordgo.Session, m *discordgo.MessageCreate, args string) error {
    return nil
}
