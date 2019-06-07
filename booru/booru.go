package booru

import (
	"strings"
)

type BooruPost struct {
	ID int
	ImageWidth int
	ImageHeight int
	PreviewFileUrl string
	FileUrl string
}

func (self *BooruPost) GetPreviewUrl() string {
	return self.PreviewFileUrl
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
