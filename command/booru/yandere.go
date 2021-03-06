package booru

import (
    "fmt"
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

func (self *yanderePost) toBooruPost() *BooruPost {
    booru_post := BooruPost {
        Source: "Yandere",
        ID: self.ID,
        URL: fmt.Sprintf("https://yande.re/post/show/%d", self.ID),
        PreviewFileUrl: self.PreviewFileUrl,
        FileUrl: self.FileUrl,
        ImageWidth: self.ImageWidth,
        ImageHeight: self.ImageHeight,
    }
    return &booru_post
}

func YandereLatestPost() (*BooruPost, error) {
    const api_url = "https://yande.re/post.json?limit=1"

    spaceClient := http.Client{Timeout: time.Second * 10}
    resp, err := spaceClient.Get(api_url)

    if err != nil {
        return nil, err
    }

    json_content, readErr := ioutil.ReadAll(resp.Body)
    if readErr != nil {
        return nil, readErr
    }

    var posts [1]yanderePost
    jsonErr := json.Unmarshal(json_content, &posts)
    if jsonErr != nil {
        return nil, jsonErr
    }

    return posts[0].toBooruPost(), nil
}
