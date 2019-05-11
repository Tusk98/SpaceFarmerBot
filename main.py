import sys
import asyncio

import discord
import toml
import xdg
from discord.ext import commands


PROGRAM_NAME = "SpaceFarmerBot"
CONFIG_FILE = "config.toml"

# initialize client
client = commands.Bot(command_prefix="!")

@client.event
async def on_ready():
    print('Logged in as')
    print(client.user.name)
    print(client.user.id)
    print('------')

if (__name__ == "__main__"):
    config_dir = xdg.XDG_CONFIG_HOME/PROGRAM_NAME
    config_file = config_dir/CONFIG_FILE
    if (config_file.exists()):
        parsed_toml = {}
        with open(config_file) as f:
            parsed_toml = toml.load(f)

        if ("token" in parsed_toml):
            client.run(parsed_toml["token"])
        else:
            print("No bot token configured!")
    else:
        print("No configuration found!")

