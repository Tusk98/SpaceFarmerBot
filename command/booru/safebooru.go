package booru

import (
    "fmt"
    "encoding/json"
    "io/ioutil"
    "net/http"
    "strings"
    "time"
)

type safebooruPost struct {
    ID int `json:"id"`
    ImageWidth int `json:"width"`
    ImageHeight int `json:"height"`
    Directory string `json:"directory"`
    Hash string `json:"hash"`
    Image string `json:"image"`
    Sample bool `json:"sample"`
}

func (self *safebooruPost) toBooruPost() *BooruPost {
    FileUrl := fmt.Sprintf("https://safebooru.org/images/%s/%s", self.Directory, self.Image)
    var PreviewFileUrl string
    if self.Sample {
        slice_ind := strings.Index(self.Image, ".")
        image_name := self.Image[:slice_ind]
        PreviewFileUrl = fmt.Sprintf("https://safebooru.org/samples/%s/sample_%s.jpg", self.Directory, image_name)
    } else {
        PreviewFileUrl = FileUrl
    }

    booru_post := BooruPost {
        Source: "Safebooru",
        ID: self.ID,
        URL: fmt.Sprintf("https://safebooru.org/index.php?page=post&s=view&id=%d", self.ID),
        PreviewFileUrl: PreviewFileUrl,
        FileUrl: FileUrl,
        ImageWidth: self.ImageWidth,
        ImageHeight: self.ImageHeight,
    }
    return &booru_post
}

func SafebooruLatestPost() (*BooruPost, error) {
    const api_url = "https://safebooru.org/index.php?page=dapi&s=post&q=index&json=1&limit=1"

    spaceClient := http.Client{Timeout: time.Second * 10}
    resp, err := spaceClient.Get(api_url)

    if err != nil {
        return nil, err
    }

    json_content, readErr := ioutil.ReadAll(resp.Body)
    if readErr != nil {
        return nil, readErr
    }

    var posts [1]safebooruPost
    jsonErr := json.Unmarshal(json_content, &posts)
    if jsonErr != nil {
        return nil, jsonErr
    }

    return posts[0].toBooruPost(), nil
}
