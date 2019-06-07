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
	"github.com/Tusk98/SpaceFarmerBot/booru"
)

const CONFIG_PATH string = "SpaceFarmerBot/config.toml"
const COMMAND_PREFIX string = "!"

type BotConfig struct {
	Token string
}

type Config struct {
	Bot BotConfig
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
	if !strings.HasPrefix(m.Content, COMMAND_PREFIX) {
		return
	}
    if strings.HasPrefix(m.Content, "!daily") {
        post, err := booru.DanbooruLatestPost()
        if err != nil {
            s.ChannelMessageSend(m.ChannelID, post.PreviewFileUrl)
        } else {
            s.ChannelMessageSend(m.ChannelID, err.Error())
        }
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
