# Discord Trivia Bot

Discord Trivia Bot is an app used to run trivia through a discord bot.
The application is made of three pieces, a UI to upload trivia through, an API to store and retrieve trivia from a Postgres database, and Javascript to run the discord bots.

## Setting up the Discord Server

In order to use this project you will need a Discord server and at least one Bot. Additionally you will need an audio bot for however many channels you'd like to have a trivia team in. Keep in mind there are various ways you could set up the bot permissions, I'm simply going tos hare how I set mine up.

### Discord Server

I set up two channels, one for bot commands and one for the trivia questions to be read out.

### Primary trivia bot

The primary trivia bot needs to be able to read and write messages in the bot commands channel and be able to write to the trivia channel. I recommend only letting this bot read the bot commands channel to limit the places it can get commands from, but that is up to you.

### Setting up an audio bot

The audio bots need to be able to join the channel which they are assigned to and the ability to use voice. I also recommend making them a priority speaker.


## Installation


### Environment Variables

You can find an example.env file in the root folder. Rename the file to '.env' and replace the placeholders with your environment variables. <br>

You will also find a .env in the tiriva-bots folder with an example file which you will need to setup. 

#### How to Run

```
docker-compose up 
```

## How to use the bots

The trivia bot reads commands sent in the bot commands channel, or any channel the bot can see, and it has various actions based on the commands. All commands are prefixed with a '!'. <br>

The main command you need is: `!start`. <br>

This will pull an unused trivia from the trivia-server, send out the image round, and then after a delay proceed to go through each round sending out the questions, and at the end will trigger the audio round, and then mark the trivia as used.  <br>

You will have to manually trigger the answer sheet with the `!answer` command. <br>


For a list of other available commands use the  `!help` command, the other commands are primarily used for triggering each part of the trivia manually.


## Configurations

There currently aren't a lot of configuration options but more could be added in the future. <br>

All configuration can be done in the config.json file located in the trivia-bots folder, the current configurations are:

- questionDelaySeconds - This sets how many seconds you want in between each question.
- imageRoundDelaySeconds - This sets how many seconds between the image round being sent out and regular trivia beginning.



## Where To Get The Trivia

Currently this project is specifically set up to handle and parse trivia from [Trivia Mafia](https://www.patreon.com/m/4171258/posts). They typically put out two full sets of trivia per week to their Patreon subscribers. <br>
https://www.patreon.com/m/4171258/posts

## Future plans

I'm currently working through various improvements and fixes as this project is still in fairly early stages and a lot of pieces have been thrown together somewhat ad hoc.<br>

I'd be open to expanding the project to work with trivia from other sources if there is interest. Feel free to reach out to me if this is something you'd be interested in. <br>



## Support
If you find any bugs or experience any issues with the project feel free to open a ticket here on github.