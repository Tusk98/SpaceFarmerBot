package booru

import (
	"fmt"
	"strings"
	"github.com/bwmarrin/discordgo"
)

type BooruPost struct {
	Source string
	ID int
	ImageWidth int
	ImageHeight int
	PreviewFileUrl string
	FileUrl string
}

func (self *BooruPost) GetPreviewUrl() string {
	return self.PreviewFileUrl
}

const COLOR int = 0xff93ac

func (self *BooruPost) ToDiscordEmbed() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed {
		Title: fmt.Sprintf("%s: #%d", self.Source, self.ID),
		Color: COLOR,
		Image: &discordgo.MessageEmbedImage {
			URL: self.PreviewFileUrl,
		},
	}
}

type UnknownBooruError struct {
	arg string
}

func (self *UnknownBooruError) Error() string {
	return self.arg
}


func BooruGetLatest(args string) (BooruPost, error) {
	if strings.HasPrefix(args, "danbooru") {
		return DanbooruLatestPost()
	} else if strings.HasPrefix(args, "yandere") {
		return YandereLatestPost()
	} else {
		return BooruPost{}, &UnknownBooruError { arg: args }
	}
}
