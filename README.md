# Standup Meeting Discord Bot
This bot is only meant to be used in my CSCE 4901 group discord server. 
Its purpose is to automate our daily standup meetings, which will be placed in separate message threads.

## Usage
Using the command `!init` in any text channel of the discord server will initialize the bot to post the daily messages in each of the specific _hard coded_ text channels.

**Note:** This `!init` command **must** be used at the time of day you want the messages to be sent. The bot simply starts a 24 hour timer once this command is run, it is not set up to automatically run at a specific time of day.
