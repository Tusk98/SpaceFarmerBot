package event

import (
    "bytes"
    "fmt"
    "strconv"
    "strings"
    "sync"
    "github.com/bwmarrin/discordgo"
)

type EventError struct {
    reason string
}
func (self *EventError) Error() string {
    return self.reason
}

type EventToml struct {
    Name string
    Date string
    Location string
    Description string
}

func (self *EventToml) toEvent() *Event {
    name := strings.TrimSpace(self.Name)
    if len(name) == 0 {
        name = "No event name specified"
    }
    date := strings.TrimSpace(self.Date)
    if len(date) == 0 {
        date = "No date specified"
    }
    location := strings.TrimSpace(self.Location)
    if len(location) == 0 {
        location = "No location specified"
    }
    description := strings.TrimSpace(self.Description)
    if len(description) == 0 {
        location = "No description available"
    }
    return &Event {
        Name: name,
        Date: date,
        Location: location,
        Description: description,
        Going: make(map[string]*discordgo.User),
    }
}

type Event struct {
    Name string
    Date string
    Location string
    Description string
    Going map[string]*discordgo.User
}

func (self *Event) ToDiscordEmbed() *discordgo.MessageEmbed {
    going := bytes.Buffer {}
    if len(self.Going) == 0 {
        going.WriteRune('-')
    } else {
        i := 1
        for _, user := range self.Going {
            going.WriteString(strconv.Itoa(i))
            going.WriteString(user.Mention())
            going.WriteRune('\n')
            i++
        }
    }
    embed := &discordgo.MessageEmbed {
        Title: self.Name,
        Color: COLOR,
        Description: self.Description,
        Fields: []*discordgo.MessageEmbedField {
            { Name: "Date", Value: self.Date },
            { Name: "Location", Value: self.Location },
            { Name: "Description", Value: self.Description },
            { Name: "Going", Value: going.String() },
        },
    }
    return embed
}

func (self *Event) ToDiscordEmbedWithID(id int) *discordgo.MessageEmbed {
    going := bytes.Buffer {}
    if len(self.Going) == 0 {
        going.WriteRune('-')
    } else {
        i := 1
        for _, user := range self.Going {
            going.WriteString(strconv.Itoa(i))
            going.WriteRune('.')
            going.WriteRune(' ')
            going.WriteString(user.Mention())
            going.WriteRune('\n')
            i++
        }
    }
    embed := &discordgo.MessageEmbed {
        Title: fmt.Sprintf("#%d %s", id, self.Name),
        Color: COLOR,
        Description: self.Description,
        Fields: []*discordgo.MessageEmbedField {
            { Name: "Date", Value: self.Date },
            { Name: "Location", Value: self.Location },
            { Name: "Description", Value: self.Description },
            { Name: "Going", Value: going.String() },
        },
    }
    return embed
}

type EventList struct {
    List []*Event
	mux sync.Mutex
}

func (self *EventList) AddEvent(event *Event) {
    self.mux.Lock()
    self.List = append(self.List, event)
    self.mux.Unlock()
}

func (self *EventList) RemoveEvent(index int) (*Event, error) {
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
        err = &EventError { reason: fmt.Sprintf("Event with ID: %d does not exist", index) }
    }
    self.mux.Unlock()

    return event, err
}

func (self *EventList) AddUserToEvent(index int, user *discordgo.User) (*Event, error) {
    var event *Event
    var err error

    self.mux.Lock()
    if index < len(self.List) {
        event = self.List[index]
        err = nil
        if _, ok := event.Going[user.ID]; !ok {
            event.Going[user.ID] = user
        } else {
            err = &EventError { reason: "You are already going to this event" }
        }
    } else {
        event = nil
        err = &EventError { reason: fmt.Sprintf("Event with ID: %d does not exist", index) }
    }
    self.mux.Unlock()

    return event, err
}

func (self *EventList) RemoveUserFromEvent(index int, userID string) (*Event, error) {
    var event *Event
    var err error

    self.mux.Lock()
    if index < len(self.List) {
        event := self.List[index]
        err = nil
        if _, ok := event.Going[userID]; ok {
            delete(event.Going, userID)
        } else {
            err = &EventError { reason: "You are already not going to this event" }
        }
    } else {
        event = nil
        err = &EventError { reason: "Event does not exist" }
    }
    self.mux.Unlock()

    return event, err
}
