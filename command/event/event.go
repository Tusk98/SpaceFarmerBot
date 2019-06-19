package event

import (
    "fmt"
    "strings"
    "strconv"
    "github.com/pelletier/go-toml"
    "github.com/bwmarrin/discordgo"
)

const COMMAND string = "event"
const DESCRIPTION string = "creates and tracks events"

const COLOR int = 0xff93ac

var eventList EventList = EventList {}

func HelpMessage(s *discordgo.Session, m *discordgo.MessageCreate) error {
    embed := &discordgo.MessageEmbed {
        Title: "event usage",
        Color: COLOR,
        Description: fmt.Sprintf("Usage: event [OPTIONS]\n%s", DESCRIPTION),
        Fields: []*discordgo.MessageEmbedField{
            { Name: "new",
                    Value: "creates a new event from values enclosed in code snippet```name = \"...\"\ndate = \"...\"\nlocation = \"...\"\ndescription = \"...\"```" },
            { Name: "remove id", Value: "removes the given event" },
            { Name: "join id", Value: "join a given event" },
            { Name: "leave id", Value: "leave a given event" },
        },
    }
    s.ChannelMessageSendEmbed(m.ChannelID, embed)
    return nil
}

func getSliceInd(args string) int {
    delims := [...]rune{ '`', ' ', '\t', '\n' }

    slice_at := -1
    for _, delim := range delims {
        slice_ind := strings.IndexRune(args, delim)
        if slice_ind != -1 {
            if slice_at == -1 {
                slice_at = slice_ind
            } else if slice_at > slice_ind {
                slice_at = slice_ind
            }
        }
    }
    return slice_at
}

func ProcessCommand(s *discordgo.Session, m *discordgo.MessageCreate, args string) error {
    s.ChannelMessageSend(m.ChannelID, "***Warning: feature is currently experimental***")

    slice_at := getSliceInd(args)
    var cmd, xs string
    if slice_at == -1 {
        cmd = args
        xs = ""
    } else {
        cmd = strings.TrimSpace(args[:slice_at])
        xs = args[slice_at:]
    }

    switch cmd {
        case "help": return HelpMessage(s, m)
        case "list": return listEvents(s, m)
        case "new": return newEvent(s, m, xs)
        case "remove": return removeEvent(s, m, xs)
        case "join": return joinEvent(s, m, xs)
        case "leave": return leaveEvent(s, m, xs)
        default: return HelpMessage(s, m)
    }
    return nil
}

func listEvents(s *discordgo.Session, m *discordgo.MessageCreate) error {
    if len(eventList.List) == 0 {
        embed := &discordgo.MessageEmbed {
            Title: "No events planned",
            Color: COLOR,
        }
        s.ChannelMessageSendEmbed(m.ChannelID, embed)
        return nil
    }
    for i, event := range eventList.List {
        embed := event.ToDiscordEmbedWithID(i)
        s.ChannelMessageSendEmbed(m.ChannelID, embed)
    }
    return nil
}

func newEvent(s *discordgo.Session, m *discordgo.MessageCreate, args string) error {
    const DISCORD_CODE_ENCLOSURE string = "```"

    event_info := args
    event_str_start := strings.Index(event_info, DISCORD_CODE_ENCLOSURE)
    if event_str_start == -1 {
        return &EventError { "Could not find start of ```" }
    }
    event_info = event_info[event_str_start + len(DISCORD_CODE_ENCLOSURE):]

    event_str_end := strings.Index(event_info, DISCORD_CODE_ENCLOSURE)
    if event_str_end == -1 {
        return &EventError { "Could not find end of ```" }
    }
    event_info = event_info[:event_str_end]

    var eventToml EventToml
    if err := toml.Unmarshal([]byte(event_info), &eventToml); err != nil {
        return err
    }
    event := eventToml.toEvent()
    id := eventList.AddEvent(event)

    embed := event.ToDiscordEmbedWithID(id)

    s.ChannelMessageSend(m.ChannelID, "New event has been added")
    s.ChannelMessageSendEmbed(m.ChannelID, embed)
/*
    s.MessageReactionAdd(msg.ChannelID, msg.ID, "\u2705")
    s.MessageReactionAdd(msg.ChannelID, msg.ID, "\u274E")
*/
    return nil
}

func removeEvent(s *discordgo.Session, m *discordgo.MessageCreate, args string) error {
    args = strings.TrimSpace(args)
    i1, err := strconv.Atoi(args)
    if err != nil {
        return err
    }
    event, err := eventList.RemoveEvent(i1)
    if err != nil {
        return err
    }
    embed := event.ToDiscordEmbedWithID(i1)
    s.ChannelMessageSend(m.ChannelID, "The following event has been removed")
    s.ChannelMessageSendEmbed(m.ChannelID, embed)
    return nil
}

func joinEvent(s *discordgo.Session, m *discordgo.MessageCreate, args string) error {
    args = strings.TrimSpace(args)
    i1, err := strconv.Atoi(args)
    if err != nil {
        return err
    }
    event, err := eventList.AddUserToEvent(i1, m.Message.Author)
    if err != nil {
        return err
    }
    embed := event.ToDiscordEmbedWithID(i1)
    s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You have been added to event #%d", i1))
    s.ChannelMessageSendEmbed(m.ChannelID, embed)
    return nil
}

func leaveEvent(s *discordgo.Session, m *discordgo.MessageCreate, args string) error {
    args = strings.TrimSpace(args)
    i1, err := strconv.Atoi(args)
    if err != nil {
        return err
    }
    event, err := eventList.RemoveUserFromEvent(i1, m.Message.Author.ID)
    if err != nil {
        return err
    }
    embed := event.ToDiscordEmbedWithID(i1)
    s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You have been removed from event #%d", i1))
    s.ChannelMessageSendEmbed(m.ChannelID, embed)
    return nil
}
