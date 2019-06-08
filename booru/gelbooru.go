package booru

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type gelbooruPost struct {
	ID int `json:"id"`
	ImageWidth int `json:"width"`
	ImageHeight int `json:"height"`
	FileUrl string `json:"file_url"`
    Directory string `json:"directory"`
    Hash string `json:"hash"`
}

func (self *gelbooruPost) toBooruPost() BooruPost {
	booru_post := BooruPost {
		Source: "Gelbooru",
		ID: self.ID,
		URL: fmt.Sprintf("https://gelbooru.com/index.php?page=post&s=view&id=%d", self.ID),
		PreviewFileUrl: fmt.Sprintf("https://img2.gelbooru.com/samples/%s/sample_%s.jpg", self.Directory, self.Hash),
		FileUrl: self.FileUrl,
		ImageWidth: self.ImageWidth,
		ImageHeight: self.ImageHeight,
	}
	return booru_post
}

func GelbooruLatestPost() (BooruPost, error) {
	const api_url = "https://gelbooru.com/index.php?page=dapi&s=post&q=index&json=1&limit=1"

	spaceClient := http.Client{Timeout: time.Second * 2}
	req, err := http.NewRequest(http.MethodGet, api_url, nil)
	if err != nil {
		return BooruPost{}, err
	}
	req.Header.Set("User-Agent", "SpaceFarmerBot")

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		return BooruPost{}, getErr
	}

	json_content, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return BooruPost{}, readErr
	}

	var posts [1]gelbooruPost
	jsonErr := json.Unmarshal(json_content, &posts)
	if jsonErr != nil {
		return BooruPost{}, jsonErr
	}

	return posts[0].toBooruPost(), nil
}
