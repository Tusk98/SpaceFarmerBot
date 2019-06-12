# SpaceFarmerBot

<p>
<img src="https://github.com/Tusk98/SpaceFarmerBot/raw/master/spacefarmer.jpg"
     width="350">
</p>

## Dependencies
 - golang (preferably 1.12+)
   - [discord-go](https://github.com/bwmarrin/discordgo)
   - [go-toml](https://github.com/pelletier/go-toml)
   - [xdg](https://github.com/adrg/xdg)

## Configuration
See [config.toml](https://github.com/Tusk98/SpaceFarmerBot/blob/master/config/config.toml)

Place config files inside `$XDG_CONFIG_HOME/SpaceFarmerBot` (usually `$HOME/.config/SpaceFarmerBot/` on GNU/Linux).

## Installing
Please make sure your go environment is set up.
```
~ $ cd SpaceFarmerBot
~ $ go install .
```
## Running
```
~ $ go run .
~ $ SpaceFarmerBot    # or if $GOPATH/bin is in your $PATH
```

## Features
 - 8ball: ask a question and it will be answered with a yes or no
 - daily: fetches latest image from supported websites
 - sauce: provide some images and sources for them will be found
