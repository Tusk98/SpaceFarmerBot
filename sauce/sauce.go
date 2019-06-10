package sauce

import (
    "bytes"
    "fmt"
    "net/http"
    "net/url"
    "io/ioutil"
    "strings"
    "github.com/bwmarrin/discordgo"
)

type GenericBotError struct {
    reason string
}
func (self *GenericBotError) Error() string {
    return self.reason
}

type SimilarityResult struct {
    URL string
    PercentSimilar uint8
    Height uint
    Width uint
}

const Command string = "sauce"
const COLOR int = 0xff93ac

const IQDB_PATTERN string = "match</th></tr><tr><td class='image'><a href=\""
const IQDB_DIMENSION_PATTERN string = "class=\"service-icon\">"

func ProcessCommand(s *discordgo.Session, m *discordgo.MessageCreate, args string) error {
    if len(m.Attachments) == 0 && len(args) == 0 {
        s.ChannelMessageSend(m.ChannelID, "No images/links to work off of? Back to farming...")
        return nil
    }

    for _, arg := range strings.Split(args, " ") {
        if len(arg) <= 6 || !strings.HasPrefix(arg, "http") {
            continue
        }

        results, err := getSimilarResults(arg)
        if err != nil {
            return err
        }
        if len(results) == 0 {
            s.ChannelMessageSend(m.ChannelID, "No results found")
            continue
        }

        var buffer bytes.Buffer
        fields := []*discordgo.MessageEmbedField{}
        for i, result := range results {
            field := &discordgo.MessageEmbedField{
                Name: fmt.Sprintf("%d: Similarity: %d%%", i+1, result.PercentSimilar),
                Value: result.URL,
            }
            fields = append(fields, field)
        }

        embed := &discordgo.MessageEmbed {
                Title: "Results",
                Color: COLOR,
                Description: buffer.String(),
                Thumbnail: &discordgo.MessageEmbedThumbnail{
                    URL: arg,
                },
                Fields: fields,
            }
        s.ChannelMessageSendEmbed(m.ChannelID, embed)
    }

    for _, attachment := range m.Attachments {
        results, err := getSimilarResults(attachment.URL)
        if err != nil {
            return err
        }

        var buffer bytes.Buffer
        for _, result := range results {
            buffer.WriteString(result.URL)
            buffer.WriteRune('\n')
        }
        embed := &discordgo.MessageEmbed {
                Title: "Results",
                Color: COLOR,
                Description: buffer.String(),
            }
        s.ChannelMessageSendEmbed(m.ChannelID, embed)
    }
    return nil
}

func getSimilarResults(img_url string) ([]SimilarityResult, error) {
    resp, err := http.PostForm("http://iqdb.org",
                url.Values{"url": { img_url }})
    if err != nil {
        return nil, err
    }
    if resp.StatusCode != 200 {
        return nil, &GenericBotError{
            reason: fmt.Sprintf("Request failed with error code: %d", resp.StatusCode),
        }
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    web_html := string(body)

    return parseResultHTML(web_html), nil
}

func parseResultHTML(body string) []SimilarityResult {
    results := []SimilarityResult{}

    yet_to_parse := body
    for i := strings.Index(yet_to_parse, IQDB_PATTERN); i != -1; i = strings.Index(yet_to_parse, IQDB_PATTERN) {
        yet_to_parse = yet_to_parse[i + len(IQDB_PATTERN):]
        idx := strings.IndexRune(yet_to_parse, '"')
        if idx == -1 {
            break
        }
        url := yet_to_parse[:idx]

        var percent uint8
        var width uint
        var height uint

        const IQDB_XY_PATTERN_1 string = "class=\"service-icon\">"
        idx = strings.Index(yet_to_parse, IQDB_XY_PATTERN_1)
        if idx == -1 {
            continue
        }
        yet_to_parse = yet_to_parse[idx + len(IQDB_XY_PATTERN_1):]

        const IQDB_XY_PATTERN_2 string = "<td>"
        idx = strings.Index(yet_to_parse, IQDB_XY_PATTERN_2)
        if idx == -1 {
            continue
        }
        yet_to_parse = yet_to_parse[idx + len(IQDB_XY_PATTERN_2):]
        fmt.Sscanf(yet_to_parse, "%d×%d", &width, &height)

        const IQDB_SIM_PATTERN_1 string = "<td>"
        idx = strings.Index(yet_to_parse, IQDB_SIM_PATTERN_1)
        if idx == -1 {
            continue
        }
        yet_to_parse = yet_to_parse[idx + len(IQDB_SIM_PATTERN_1):]

        fmt.Sscanf(yet_to_parse, "%d%% similarity", &percent)

        // some urls might not be complete
        if !strings.HasPrefix(url, "http") {
            url = fmt.Sprintf("https:%s", url)
        }

        result := SimilarityResult {
            URL: url,
            Width: width,
            Height: height,
            PercentSimilar: percent,
        }
        fmt.Printf("%+v\n", result)
        results = append(results, result)
    }
    return results
}
