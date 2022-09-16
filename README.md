# Standup Meeting Discord Bot
This bot is only meant to be used in my CSCE 4901 group discord server. 
Its purpose is to automate our daily standup meetings, which will be placed in separate message threads.

## Usage
Just start the bot and you're ready to go, no initialization command needed. The bot will find the standup channels automatically based on their parent category (_hard coded ID_). It will then calculate the remaining time until 8:00 AM and wait to send messages until then. 

Responses sent within the individual standup threads are saved to a text file for safe keeping, should they be needed for documentation at a later time. Using `!getResponses <your_standup_channel_name>` will upload the corresponding text file for your use.
