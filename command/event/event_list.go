package event

import (
    "fmt"
    "strconv"
    "strings"
    "sync"
    "github.com/bwmarrin/discordgo"
    "github.com/pelletier/go-toml"
    "github.com/Tusk98/SpaceFarmerBot/command"
)

const COMMAND string = "event"
const DESCRIPTION string = "creates, manages and tracks events"

const COLOR int = 0xff93ac


type EventList struct {
    List []*Event
	mux sync.Mutex
}
/* for BotCommand interface */
func (self *EventList) Prefix() string {
    return COMMAND
}
/* for BotCommand interface */
func (self *EventList) Description() string {
    return DESCRIPTION
}
/* for BotCommand interface */
func (self *EventList) HelpMessage(s *discordgo.Session, m *discordgo.MessageCreate) error {
    embed := &discordgo.MessageEmbed {
        Title: "event usage",
        Color: COLOR,
        Description: fmt.Sprintf("Usage: event [OPTIONS]\n%s", DESCRIPTION),
        Fields: []*discordgo.MessageEmbedField{
            { Name: "new",
                    Value: "creates a new event from values enclosed in code snippet" +
                            "```name = \"...\"\ndate = \"...\"\nlocation = \"...\"\ndescription = \"...\"```" },
            { Name: "remove id", Value: "removes the given event" },
            { Name: "join id", Value: "join a given event" },
            { Name: "leave id", Value: "leave a given event" },
        },
    }
    s.ChannelMessageSendEmbed(m.ChannelID, embed)
    return nil
}
/* for BotCommand interface */
func (self *EventList) ProcessCommand(s *discordgo.Session, m *discordgo.MessageCreate, args string) error {
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
        case "help": return self.HelpMessage(s, m)
        case "list": return self.listEvents(s, m)
        case "new": return self.newEvent(s, m, xs)
        case "remove": return self.removeEvent(s, m, xs)
        case "join": return self.joinEvent(s, m, xs)
        case "leave": return self.leaveEvent(s, m, xs)
        default: return self.HelpMessage(s, m)
    }
    return nil
}

/* atomically add an Event to the list of events */
func (self *EventList) AddEvent(event *Event) int {
    self.mux.Lock()
    self.List = append(self.List, event)
    index := len(self.List) - 1
    self.mux.Unlock()
    return index
}

/* atomically delete an Event from the list of events */
func (self *EventList) DeleteEvent(index int) (*Event, error) {
    var event *Event
    var err error

    self.mux.Lock()
    if index < len(self.List) {
        event = self.List[index]
        err = nil
        self.List[index] = self.List[len(self.List) - 1]
        self.List = self.List[:len(self.List) - 1]
    } else {
        event = nil
        err = &command.CommandError { Reason: fmt.Sprintf("Event with ID: %d does not exist", index) }
    }
    self.mux.Unlock()

    return event, err
}

/* atomically add a user to an Event's Going list */
func (self *EventList) EventUserAdd(index int, user *discordgo.User) (*Event, error) {
    var event *Event
    var err error

    self.mux.Lock()
    if index < len(self.List) {
        event = self.List[index]
        if _, exists := event.Going[user.ID]; exists {
            event.Going[user.ID] = user
            err = nil
        } else {
            event = nil
            err = &command.CommandError {
                Reason: "You are already going to this event",
            }
        }
    } else {
        event = nil
        err = &command.CommandError {
                Reason: fmt.Sprintf("Event with ID: %d does not exist", index),
            }
    }

    self.mux.Unlock()

    return event, err
}

/* atomically remove a user from an Event's Going list */
func (self *EventList) EventUserRemove(index int, userID string) (*Event, error) {
    var event *Event
    var err error

    self.mux.Lock()
    if index < len(self.List) {
        event = self.List[index]
        if _, exists := event.Going[userID]; exists {
            delete(event.Going, userID)
            err = nil
        } else {
            event = nil
            err = &command.CommandError {
                Reason: "You were not on the list of people going",
            }
        }
    } else {
        event = nil
        err = &command.CommandError {
                Reason: fmt.Sprintf("Event with ID: %d does not exist", index),
            }
    }
    self.mux.Unlock()

    return event, err
}

/* parses user input for AddEvent */
func (self *EventList) newEvent(s *discordgo.Session, m *discordgo.MessageCreate, args string) error {
    const DISCORD_CODE_ENCLOSURE string = "```"

    /* find toml string */
    event_info := args
    event_str_start := strings.Index(event_info, DISCORD_CODE_ENCLOSURE)
    if event_str_start == -1 {
        return &command.CommandError { Reason: "Could not find start of ```" }
    }
    event_info = event_info[event_str_start + len(DISCORD_CODE_ENCLOSURE):]

    event_str_end := strings.Index(event_info, DISCORD_CODE_ENCLOSURE)
    if event_str_end == -1 {
        return &command.CommandError { Reason: "Could not find end of ```" }
    }
    event_info = event_info[:event_str_end]

    /* parse toml string */
    var eventToml EventToml
    if err := toml.Unmarshal([]byte(event_info), &eventToml); err != nil {
        return err
    }

    /* convert toml to event, add it to list of events and send a message */
    event := eventToml.toEvent()
    id := self.AddEvent(event)

    embed := event.ToDiscordEmbedWithID(id)

    s.ChannelMessageSend(m.ChannelID, "New event has been added")
    s.ChannelMessageSendEmbed(m.ChannelID, embed)
/*
    s.MessageReactionAdd(msg.ChannelID, msg.ID, "\u2705")
    s.MessageReactionAdd(msg.ChannelID, msg.ID, "\u274E")
*/
    return nil
}

/* parses user input for DeleteEvent */
func (self *EventList) removeEvent(s *discordgo.Session, m *discordgo.MessageCreate, args string) error {
    args = strings.TrimSpace(args)
    i1, err := strconv.Atoi(args)
    if err != nil {
        return err
    }
    event, err := self.DeleteEvent(i1)
    if err != nil {
        return err
    }
    embed := event.ToDiscordEmbedWithID(i1)
    s.ChannelMessageSend(m.ChannelID, "The following event has been removed")
    s.ChannelMessageSendEmbed(m.ChannelID, embed)
    return nil
}

/* parses user input for EventUserAdd */
func (self *EventList) joinEvent(s *discordgo.Session, m *discordgo.MessageCreate, args string) error {
    args = strings.TrimSpace(args)
    i1, err := strconv.Atoi(args)
    if err != nil {
        return err
    }
    event, err := self.EventUserAdd(i1, m.Message.Author)
    if err != nil {
        return err
    }
    embed := event.ToDiscordEmbedWithID(i1)
    s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You have been added to event #%d", i1))
    s.ChannelMessageSendEmbed(m.ChannelID, embed)
    return nil
}

/* parses user input for EventUserDelete */
func (self *EventList) leaveEvent(s *discordgo.Session, m *discordgo.MessageCreate, args string) error {
    args = strings.TrimSpace(args)
    i1, err := strconv.Atoi(args)
    if err != nil {
        return err
    }
    event, err := self.EventUserRemove(i1, m.Message.Author.ID)
    if err != nil {
        return err
    }
    embed := event.ToDiscordEmbedWithID(i1)
    s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You have been removed from event #%d", i1))
    s.ChannelMessageSendEmbed(m.ChannelID, embed)
    return nil
}

func (self *EventList) listEvents(s *discordgo.Session, m *discordgo.MessageCreate) error {
    if len(self.List) == 0 {
        embed := &discordgo.MessageEmbed {
            Title: "No events planned",
            Color: COLOR,
        }
        s.ChannelMessageSendEmbed(m.ChannelID, embed)
        return nil
    }
    for i, event := range self.List {
        embed := event.ToDiscordEmbedWithID(i)
        s.ChannelMessageSendEmbed(m.ChannelID, embed)
    }
    return nil
}

/* ugly way of splitting into subcommand */
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
