package event

import (
    "bytes"
    "fmt"
    "strings"
    "strconv"
    "github.com/bwmarrin/discordgo"
)

type EventToml struct {
    Name string
    Date string
    Location string
    Description string
}

/* produces an Event from a EventToml received from user input */
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

/* used to keep track of details of an Event */
type Event struct {
    Name string
    Date string
    Location string
    Description string
    Going map[string]*discordgo.User
}

/* Produces a discord embed message from an Event */
func (self *Event) ToDiscordEmbedWithID(id int) *discordgo.MessageEmbed {
    going := bytes.Buffer {}
    if len(self.Going) == 0 {
        going.WriteRune('-')
    } else {
        i := 1
        for _, user := range self.Going {
            going.WriteString(strconv.Itoa(i))
            going.WriteString(". ")
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
