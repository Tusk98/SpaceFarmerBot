
all: SpaceFarmerBot

SpaceFarmerBot: main.go
	go install github.com/tusk98/SpaceFarmerBot

dependencies:
	go get github.com/adrg/xdg
	go get github.com/bwmarrin/discordgo
	go get github.com/pelletier/go-toml
