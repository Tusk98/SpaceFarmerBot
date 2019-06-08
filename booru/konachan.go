package booru

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type konachanPost struct {
	ID int `json:"id"`
	ImageWidth int `json:"width"`
	ImageHeight int `json:"height"`
	PreviewFileUrl string `json:"sample_url"`
	FileUrl string `json:"file_url"`
}

func (self *konachanPost) toBooruPost() BooruPost {
	booru_post := BooruPost {
		Source: "Konachan",
		ID: self.ID,
		URL: fmt.Sprintf("https://konachan.com/post/show/%d", self.ID),
		PreviewFileUrl: self.PreviewFileUrl,
		FileUrl: self.FileUrl,
		ImageWidth: self.ImageWidth,
		ImageHeight: self.ImageHeight,
	}
	return booru_post
}

func KonachanLatestPost() (BooruPost, error) {
	const api_url = "https://konachan.com/post.json?limit=1"

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

	var posts [1]konachanPost
	jsonErr := json.Unmarshal(json_content, &posts)
	if jsonErr != nil {
		return BooruPost{}, jsonErr
	}

	return posts[0].toBooruPost(), nil
}
