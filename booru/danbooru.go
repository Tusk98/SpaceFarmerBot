package booru

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type danbooruPost struct {
	ID int `json:"id"`
	ImageWidth int `json:"image_width"`
	ImageHeight int `json:"image_height"`
	PreviewFileUrl string `json:"large_file_url"`
	FileUrl string `json:"file_url"`
}

func (self *danbooruPost) toBooruPost() BooruPost {
	booru_post := BooruPost {
		Source: "Danbooru",
		ID: self.ID,
		URL: fmt.Sprintf("https://danbooru.donmai.us/posts/%d", self.ID),
		PreviewFileUrl: self.PreviewFileUrl,
		FileUrl: self.FileUrl,
		ImageWidth: self.ImageWidth,
		ImageHeight: self.ImageHeight,
	}
	return booru_post
}

func DanbooruLatestPost() (BooruPost, error) {
	const api_url = "https://danbooru.donmai.us/posts.json?limit=1"

	spaceClient := http.Client{Timeout: time.Second * 5}
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

	var posts [1]danbooruPost
	jsonErr := json.Unmarshal(json_content, &posts)
	if jsonErr != nil {
		return BooruPost{}, jsonErr
	}

	return posts[0].toBooruPost(), nil
}
