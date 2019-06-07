package booru

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type yanderePost struct {
	ID int `json:"id"`
	ImageWidth int `json:"width"`
	ImageHeight int `json:"height"`
	PreviewFileUrl string `json:"sample_url"`
	FileUrl string `json:"file_url"`
}

func (self *yanderePost) toBooruPost() BooruPost {
	booru_post := BooruPost {
		Source: "Yandere",
		ID: self.ID,
		ImageWidth: self.ImageWidth,
		ImageHeight: self.ImageHeight,
		PreviewFileUrl: self.PreviewFileUrl,
		FileUrl: self.FileUrl,
	}
	return booru_post
}

func YandereLatestPost() (BooruPost, error) {
	const api_url = "https://yande.re/post.json?limit=1"

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

	var posts [1]yanderePost
	jsonErr := json.Unmarshal(json_content, &posts)
	if jsonErr != nil {
		return BooruPost{}, jsonErr
	}

	return posts[0].toBooruPost(), nil
}
