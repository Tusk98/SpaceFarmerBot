package foundation

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/Tusk98/SpaceFarmerBot/command"
	"github.com/bwmarrin/discordgo"
	"github.com/pelletier/go-toml"
)

const COMMAND string = "foundation"
const DESCRIPTION string = "SCP Foundation section of the Bot"

const COLOR int = 0xff93ac
