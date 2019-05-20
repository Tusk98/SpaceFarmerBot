package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"github.com/adrg/xdg"
	"github.com/bwmarrin/discordgo"
	"github.com/pelletier/go-toml"
)

type BotConfig struct {
	Token string
}

type Config struct {
	Bot BotConfig
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}
const CONFIG_PATH string = "SpaceFarmerBot/config.toml"

func main() {
	config_path, err := xdg.SearchConfigFile(CONFIG_PATH)
	if err != nil {
		fmt.Println("Failed to find config file:", err)
		return
	}
	toml_data, err := ioutil.ReadFile(config_path)
	if err != nil {
		fmt.Println("Failed to read config file:", err)
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
	dg.AddHandler(messageCreate)

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
