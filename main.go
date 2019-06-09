package main

import (
    "fmt"
    "io/ioutil"
    "strings"
    "os"
    "os/signal"
    "syscall"
    "github.com/adrg/xdg"
    "github.com/bwmarrin/discordgo"
    "github.com/pelletier/go-toml"
    "github.com/Tusk98/SpaceFarmerBot/ball"
    "github.com/Tusk98/SpaceFarmerBot/booru"
    "github.com/Tusk98/SpaceFarmerBot/sauce"
)

const CONFIG_PATH string = "SpaceFarmerBot/config.toml"
const COMMAND_PREFIX string = "!"
const COLOR int = 0xff93ac

type BotConfig struct {
    Token string
}

type Config struct {
    Bot BotConfig
}

type UnknownCommandError struct {
    arg string
}
func (self *UnknownCommandError) Error() string {
    return self.arg
}


func onReady(discord *discordgo.Session, ready *discordgo.Ready) {
    err := discord.UpdateStatus(0, "A friendly helpful bot!")
    if err != nil {
        fmt.Println("Error attempting to set bot status:", err)
    }
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func commandHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
    // do not do anything if message is from bot
    if m.Author.ID == s.State.User.ID {
        return
    }
    // do nothing if message does not start with command invoking string
    if !strings.HasPrefix(m.Content, COMMAND_PREFIX) {
        return
    }

    msg := m.Content[len(COMMAND_PREFIX):]
    slice_ind := strings.IndexRune(m.Content, ' ')

    /* sliced as a space
     * e.g. "!8ball   answer my question " becomes:
     *    command = "8ball"
     *    args = "answer my question"
     */
    var command, args string
    if slice_ind != -1 {
        command = msg[:slice_ind-1]
        args = strings.TrimSpace(msg[slice_ind:])
    } else {
        command = msg
        args = ""
    }
    fmt.Printf("cmd: \"%s\"\nargs: \"%s\"\n", command, args)

    var err error = nil

    // check for valid commands
    switch command {
    case booru.Command: err = booru.ProcessCommand(s, m, args)
    case ball.Command: err = ball.ProcessCommand(s, m, args)
    case sauce.Command: err = sauce.ProcessCommand(s, m, args)
    // unknown command
    default: err = &UnknownCommandError {
            arg: fmt.Sprintf("Unknown Command: %s", command),
        }
    }

    if err != nil {
        s.ChannelMessageSend(m.ChannelID, err.Error())
    }
}

func main() {
    config_path, err := xdg.SearchConfigFile(CONFIG_PATH)
    if err != nil {
        fmt.Println("Failed to find config file:", err)
        return
    }
    toml_data, err := ioutil.ReadFile(config_path)
    if err != nil {
        fmt.Println("Failed to read config file:", err)
        return
    }

    var config Config
    if err := toml.Unmarshal(toml_data, &config); err != nil {
        fmt.Println("Failed to parse config file:", err)
        return
    }
    var Token = config.Bot.Token
    fmt.Println("Token", Token)

    // Create a new Discord session using the provided bot token.
    dg, err := discordgo.New("Bot " + Token)
    if err != nil {
        fmt.Println("error creating Discord session,", err)
        return
    }

    // Register the messageCreate func as a callback for MessageCreate events.
    dg.AddHandler(onReady)
    dg.AddHandler(commandHandler)

    // Open a websocket connection to Discord and begin listening.
    err = dg.Open()
    if err != nil {
        fmt.Println("error opening connection,", err)
        return
    }

    // Wait here until CTRL-C or other term signal is received.
    fmt.Println("Bot is now running.  Press CTRL-C to exit.")
    sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

    <-sc

    // Cleanly close down the Discord session.
    dg.Close()
}
