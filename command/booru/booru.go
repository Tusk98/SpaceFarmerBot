package booru

import (
    "fmt"
    "math/rand"
    "github.com/bwmarrin/discordgo"
    "github.com/Tusk98/SpaceFarmerBot/command"
)

const COMMAND string = "daily"
const DESCRIPTION string = "fetches latest image from supported websites"

const COLOR int = 0xff93ac
const BOORUS_SUPPORTED uint = 5
const (
    Danbooru    = iota
    Gelbooru    = iota
    Konachan    = iota
    Safebooru   = iota
    Yandere     = iota
)

type BooruPost struct {
    Source string
    ID int
    URL string
    PreviewFileUrl string
    FileUrl string
    ImageWidth int
    ImageHeight int
}
func (self *BooruPost) ToDiscordEmbed() *discordgo.MessageEmbed {
    return &discordgo.MessageEmbed {
        Title: fmt.Sprintf("%s: #%d", self.Source, self.ID),
        Color: COLOR,
        Image: &discordgo.MessageEmbedImage {
            URL: self.PreviewFileUrl,
        },
        URL: self.URL,
    }
}

type BooruCommand struct {}

/* for BotCommand interface */
func (self *BooruCommand) Prefix() string {
    return COMMAND
}
/* for BotCommand interface */
func (self *BooruCommand) Description() string {
    return DESCRIPTION
}
/* for BotCommand interface */
func (self *BooruCommand) HelpMessage(s *discordgo.Session, m *discordgo.MessageCreate) error {
    embed := &discordgo.MessageEmbed {
        Title: "daily usage",
        Color: COLOR,
        Description: fmt.Sprintf("Usage: daily [OPTIONS]\n%s", DESCRIPTION),
        Fields: []*discordgo.MessageEmbedField {
            { Name: "all", Value: "fetches all the latest images from supported platforms" },
            { Name: "danbooru", Value: "fetches the latest image from danbooru" },
            { Name: "gelbooru", Value: "fetches latest image on gelbooru" },
            { Name: "konachan", Value: "fetches latest image on konachan" },
            { Name: "safebooru", Value: "fetches latest image on safebooru" },
            { Name: "yandere", Value: "fetches latest image on yandere" },
        },
    }
    s.ChannelMessageSendEmbed(m.ChannelID, embed)
    return nil
}
/* for BotCommand interface */
func (self *BooruCommand) ProcessCommand(s *discordgo.Session, m *discordgo.MessageCreate, args string) error {
    switch args {
        case "help": return self.HelpMessage(s, m)
        case "": {
            booru := uint(rand.Int()) % BOORUS_SUPPORTED
            post, err := BooruGetLatest(booru)
            if err != nil {
                return err
            }
            embed := post.ToDiscordEmbed()
            s.ChannelMessageSendEmbed(m.ChannelID, embed)
        }
        case "all": {
            for i := 0; uint(i) < BOORUS_SUPPORTED; i++ {
                post, err := BooruGetLatest(uint(i))
                if err != nil {
                    return err
                }
                embed := post.ToDiscordEmbed()
                s.ChannelMessageSendEmbed(m.ChannelID, embed)
            }
        }
        default: {
            booru, err := parseBooruType(args)
            if err != nil {
                return err
            }
            post, err := BooruGetLatest(booru)
            if err != nil {
                return err
            }
            embed := post.ToDiscordEmbed()
            s.ChannelMessageSendEmbed(m.ChannelID, embed)
        }
    }
    return nil
}

func parseBooruType(arg string) (uint, error) {
    switch arg {
    case "danbooru": return Danbooru, nil
    case "yandere": return Yandere, nil
    case "konachan": return Konachan, nil
    case "safebooru": return Safebooru, nil
    case "gelbooru": return Gelbooru, nil
    default: return 0, &command.CommandError { Reason: fmt.Sprintf("Unknown argument: %s", arg) }
    }
}

func BooruGetLatest(booru uint) (*BooruPost, error) {
    switch booru {
    case Danbooru: return DanbooruLatestPost()
    case Yandere: return YandereLatestPost()
    case Konachan: return KonachanLatestPost()
    case Safebooru: return SafebooruLatestPost()
    case Gelbooru: return GelbooruLatestPost()
    default: return nil, &command.CommandError {}
    }
}
