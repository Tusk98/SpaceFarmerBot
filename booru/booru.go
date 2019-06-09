package booru

import (
    "fmt"
    "math/rand"
    "github.com/bwmarrin/discordgo"
)

const Command string = "daily"


const COLOR int = 0xff93ac
const BOORUS_SUPPORTED uint = 5
const (
    Danbooru    = iota
    Gelbooru    = iota
    Konachan    = iota
    Safebooru    = iota
    Yandere        = iota
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

type UnknownBooruError struct {
    arg string
}
func (self *UnknownBooruError) Error() string {
    return self.arg
}

func ProcessCommand(s *discordgo.Session, m *discordgo.MessageCreate, args string) error {
    if args == "" {
        booru := uint(rand.Int()) % BOORUS_SUPPORTED
        post, err := BooruGetLatest(booru)
        if err != nil {
            return err
        }
        embed := post.ToDiscordEmbed()
        s.ChannelMessageSendEmbed(m.ChannelID, embed)
    } else if args == "all" {
        for i := 0; uint(i) < BOORUS_SUPPORTED; i++ {
            post, err := BooruGetLatest(uint(i))
            if err != nil {
                return err
            }
            embed := post.ToDiscordEmbed()
            s.ChannelMessageSendEmbed(m.ChannelID, embed)
        }
    } else {
        booru, err := parseBooruType(args)
        if err != nil {
            return err
        }
        post, err := BooruGetLatest(booru)
        if err != nil {
            return err
        }
        fmt.Printf("%+v\n", post)
        embed := post.ToDiscordEmbed()
        s.ChannelMessageSendEmbed(m.ChannelID, embed)
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
    default: return 100, &UnknownBooruError { arg: fmt.Sprintf("Unknown argument: %s", arg) }
    }
}

func BooruGetLatest(booru uint) (BooruPost, error) {
    switch booru {
    case Danbooru: return DanbooruLatestPost()
    case Yandere: return YandereLatestPost()
    case Konachan: return KonachanLatestPost()
    case Safebooru: return SafebooruLatestPost()
    case Gelbooru: return GelbooruLatestPost()
    default: return BooruPost{}, &UnknownBooruError {}
    }
}
