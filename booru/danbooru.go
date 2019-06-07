package booru

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type DanbooruPost struct {
    ID int `json:"id"`
    ImageWidth int `json:"image_width"`
    ImageHeight int `json:"image_height"`
	PreviewFileUrl string `json:"preview_file_url"`
    FileUrl string `json:"file_url"`
    TagsGeneral string `json:"tag_string_general"`
    TagsCharacter string `json:"tag_string_character"`
    TagsCopyright string `json:"tag_string_copyright"`
    TagsArtist string `json:"tag_string_artist"`
    TagsMeta string `json:"tag_string_meta"`
}

func (self DanbooruPost) toBooruPost() BooruPost {
    booru_post := BooruPost {
        ID: self.ID,
        ImageWidth: self.ImageWidth,
        ImageHeight: self.ImageHeight,
	    PreviewFileUrl: self.PreviewFileUrl,
        FileUrl: self.FileUrl,
        TagsGeneral: self.TagsGeneral,
        TagsCharacter: self.TagsCharacter,
        TagsCopyright: self.TagsCopyright,
        TagsArtist: self.TagsArtist,
        TagsMeta: self.TagsMeta,
    }
    return booru_post
}

func DanbooruLatestPost() (BooruPost, error) {
    api_url := "https://danbooru.donmai.us/posts.json?limit=1"

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

	var danbooru_posts [1]DanbooruPost
	jsonErr := json.Unmarshal(json_content, &danbooru_posts)
	if jsonErr != nil {
        return BooruPost{}, jsonErr
	}

    return danbooru_posts[0].toBooruPost(), nil
}
