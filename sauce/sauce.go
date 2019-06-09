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

const Command string = "sauce"
const COLOR int = 0xff93ac

func ProcessCommand(s *discordgo.Session, m *discordgo.MessageCreate, args string) error {
    if len(m.Attachments) == 0 {
        s.ChannelMessageSend(m.ChannelID, "No images to work off of? Back to farming...")
        return nil
    }

    const pattern string = "match</th></tr><tr><td class='image'><a href=\""

    spaceClient := http.Client{}
    for _, attachment := range m.Attachments {
        resp, err := spaceClient.PostForm("http://iqdb.org",
                url.Values{"url": { attachment.URL }})
        if err != nil {
            return err
        }
        if resp.StatusCode != 200 {
            return &GenericBotError{
                reason: fmt.Sprintf("Request failed with error code: %d", resp.StatusCode),
            }
        }
        defer resp.Body.Close()

        body, err := ioutil.ReadAll(resp.Body)
        web_html := string(body)

        var buffer bytes.Buffer
        yet_to_parse_html := web_html
        for i := strings.Index(yet_to_parse_html, pattern); i != -1; i = strings.Index(yet_to_parse_html, pattern) {
            yet_to_parse_html = yet_to_parse_html[i+len(pattern):]
            idx := strings.IndexRune(yet_to_parse_html, '"')
            if idx == -1 {
                break
            }
            url := yet_to_parse_html[:idx]
            if !strings.HasPrefix(url, "http") {
                buffer.WriteString("https:")
            }
            buffer.WriteString(url)
            buffer.WriteString("\n")
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
