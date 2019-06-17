# SpaceFarmerBot

<p>
<img src="https://github.com/Tusk98/SpaceFarmerBot/raw/master/spacefarmer.jpg"
     width="350">
</p>

## Dependencies
 - golang (preferably 1.12+, older versions require you to manually download dependencies via `go get`)
   - [discord-go](https://github.com/bwmarrin/discordgo)
   - [go-toml](https://github.com/pelletier/go-toml)
   - [xdg](https://github.com/adrg/xdg)

## Configuration
See [config.toml](https://github.com/Tusk98/SpaceFarmerBot/blob/master/config/config.toml)

Place config files inside `$XDG_CONFIG_HOME/SpaceFarmerBot` (usually `$HOME/.config/SpaceFarmerBot/` on GNU/Linux).

## Building
Please make sure your go environment is set up.
```
~ $ cd SpaceFarmerBot
~ $ go build .
```

## Installing
```
~ $ go install .
```
## Running
Running without installing
```
~ $ go run .
```
If you've already installed it and `$GOPATH/bin` is in your `$PATH`
```
~ $ SpaceFarmerBot
```
## Features
 - 8ball: ask a question and it will be answered with a yes or no
```
# example usage
+8ball is it going to rain today?
+8ball are the raptors going to win?
```
 - daily: fetches latest image from supported websites
```
# example usage
+daily
+daily danbooru
+daily yandere
+daily all
```
 - event: manage upcoming events
```
# example usage
+event new ` ` ` <-- note these should not contain spaces in between
name = "New Event"
description = "This is a new event"
date = "This is a string that represents a date"
location = "This is where the event will take place"
` ` `
+event remove 0
+event join 0
+event leave 0
+event list
```
 - sauce: provide some images and sources for them will be found
```
# example usage
+sauce         # has an image attached
+sauce https://github.com/Tusk98/SpaceFarmerBot/raw/master/spacefarmer.jpg
```
