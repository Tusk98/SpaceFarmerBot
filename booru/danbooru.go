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

    spaceClient := http.Client{Timeout: time.Second * 10}
    resp, err := spaceClient.Get(api_url)

    if err != nil {
        return BooruPost{}, err
    }

    json_content, readErr := ioutil.ReadAll(resp.Body)
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
