package run

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Tusk98/SpaceFarmerBot/command"
	"github.com/Tusk98/SpaceFarmerBot/command/ball"
	"github.com/Tusk98/SpaceFarmerBot/command/booru"
	"github.com/Tusk98/SpaceFarmerBot/command/event"
	"github.com/Tusk98/SpaceFarmerBot/command/sauce"
	"github.com/Tusk98/SpaceFarmerBot/command/tldr"
	"github.com/Tusk98/SpaceFarmerBot/config"
	"github.com/bwmarrin/discordgo"
)

const COMMAND_PREFIX string = "+"
const COLOR int = 0xff93ac

var _STATUS_VALUES []string = []string{
	"Bargaining with Maroo",
	"Completing the codex",
	"Extracting Nitain",
	"Failing sortie spy",
	"Finding Kurias",
	"Headpatting noggles",
	"Helping Clem",
	"Looking for frost leaves",
	"Sabotaging Vay Hek's plans",
	"Shopping for syandanas",
	"Space farming argon crystals",
	"Unveiling rivens",
}

type UnknownCommandError struct {
	arg string
}

func (self *UnknownCommandError) Error() string {
	return self.arg
}

var cmdDict map[string]command.BotCommand = initCmdDict()

func initCmdDict() map[string]command.BotCommand {
	dict := make(map[string]command.BotCommand)

	// help command
	help := HelpCommand{}
	dict[help.Prefix()] = &help

	// 8ball command
	eight_ball := ball.EightBall{}
	dict[eight_ball.Prefix()] = &eight_ball

	// daily command
	booru := booru.BooruCommand{}
	dict[booru.Prefix()] = &booru

	// event command
	event_list := event.EventList{}
	dict[event_list.Prefix()] = &event_list

	// sauce command
	sauce := sauce.SauceCommand{}
	dict[sauce.Prefix()] = &sauce

	// sauce command
	tldr := tldr.TldrCommand{}
	dict[tldr.Prefix()] = &tldr

	foundation := foundation.FoundationCommand{}
	dict[foundation.Prefix()] = &foundation

	return dict
}

type HelpCommand struct{}

/* for BotCommand interface */
func (self *HelpCommand) Prefix() string {
	return "help"
}

/* for BotCommand interface */
func (self *HelpCommand) Description() string {
	return "provides help for all commands"
}

/* for BotCommand interface */
func (self *HelpCommand) HelpMessage(s *discordgo.Session, m *discordgo.MessageCreate) error {
	fields := []*discordgo.MessageEmbedField{}
	for key, cmd := range cmdDict {
		field := discordgo.MessageEmbedField{Name: key, Value: cmd.Description()}
		fields = append(fields, &field)
	}

	embed := &discordgo.MessageEmbed{
		Title:       "SpaceFarmerBot Usage",
		Color:       COLOR,
		Description: fmt.Sprintf("%scommand arguments", COMMAND_PREFIX),
		Fields:      fields,
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return nil
}

/* for BotCommand interface */
func (self *HelpCommand) ProcessCommand(s *discordgo.Session, m *discordgo.MessageCreate, args string) error {
	if val, ok := cmdDict[args]; ok {
		return val.HelpMessage(s, m)
	} else {
		return self.HelpMessage(s, m)
	}
}

func onReady(discord *discordgo.Session, ready *discordgo.Ready) {
	status_ind := rand.Int() % len(_STATUS_VALUES)
	status := _STATUS_VALUES[status_ind]
	err := discord.UpdateStatus(0, status)
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

	/* sliced as a space
	 * e.g. "+8ball   answer my question " becomes:
	 *    command = "8ball"
	 *    args = "answer my question"
	 */
	var command, args string
	slice_ind := strings.IndexRune(m.Content, ' ')
	if slice_ind != -1 {
		command = msg[:slice_ind-1]
		args = strings.TrimSpace(msg[slice_ind:])
	} else {
		command = msg
		args = ""
	}
	fmt.Printf("cmd: \"%s\"\nargs: \"%s\"\n", command, args)

	if val, ok := cmdDict[command]; ok {
		err := val.ProcessCommand(s, m, args)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
		}
	}
}

func reactHandler(s *discordgo.Session, m *discordgo.MessageReactionAdd) {}

func exitHandler(status int) {
	fmt.Println("Exiting...")
	os.Exit(status)
	//    toml.Marshal
}

func Run() {
	config, err := config.GetConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var Token = config.Bot.Token
	fmt.Println("Token", Token)

	/* Create a new Discord session using the provided bot token. */
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session:", err)
		os.Exit(2)
	}

	/* Register functions as a callback for MessageCreate events */
	dg.AddHandler(onReady)
	dg.AddHandler(commandHandler)
	//    dg.AddHandler(reactHandler)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection:", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}
