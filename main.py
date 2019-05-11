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
    if (not config_file.exists()):
        print("No configuration found!")
        exit(1)

    modules = []
    for module in modules:
        bot.load_extension(module)

    parsed_toml = {}
    with open(config_file) as f:
        parsed_toml = toml.load(f)

    if ("bot" not in parsed_toml):
        print("No bot configuration found!")
        exit(1)

    bot_config = parsed_toml["bot"]
    if ("token" not in bot_config):
        print("No bot token configured!")
        exit(1)

    token = bot_config["token"]
    print("Logging in...")
    try:
        asyncio.get_event_loop().run_until_complete(bot.start(token))
    except KeyboardInterrupt:
        asyncio.get_event_loop().run_until_complete(tally_before_exit())
        asyncio.get_event_loop().run_until_complete(bot.logout())
        # cancel all tasks lingering

    finally:
        asyncio.get_event_loop().run_until_complete(bot.close())
        asyncio.get_event_loop().run_until_complete(asyncio.gather(*asyncio.Task.all_tasks()))
