# Trivia Bots

This piece of the project contains the code for the main trivia bot which handles reading commands and sending out questions, as well as the code for any audio buts to play audio rounds.

## Setting up the Bots

In order to use this project you will need a Discord server and at least one Bot. Additionally you will need an audio bot for however many channels you'd like to have a trivia team in. Keep in mind there are various ways you could set up the bot permissions, I'm simply going tos hare how I set mine up.

### Setting up primary trivia bot

TODO

### Setting up an audio bot

TODO


## Installation

### Prerequisites

You will need to have Node.js installed to run this project. <br>

The audio bots require FFmpeg to be installed so make sure to do so if you plan to have an audio round.<br>


### Environment Variables

You can find an example.env file in the trivia-bots folder. Rename the file to '.env' and replace the placeholders with your environment variables.

#### How to Run

Install dependencies:<br>
```
npm install
```

Start the application:<br>
```
node index.js
```


## How to use the bots

The trivia bot reads commands sent in the bot commands channel, or any channel the bot can see, and it has various actions based on the commands. All commands are prefixed with a '!'. <br>

The main command you need is: `!start`. <br>

This will pull an unused trivia from the trivia-server, send out the image round, and then after a delay proceed to go through each round sending out the questions, and at the end will trigger the audio round, and then mark the trivia as used.  <br>

You will have to manually trigger the answer sheet with the `!answer` command. <br>


For a list of other available commands use the  `!help` command, the other commands are primarily used for triggering each part of the trivia manually.


## Configurations

There currently aren't a lot of configuration options but more could be added in the future. <br>

All configuration can be done in the config.json file located in this folder, the current configurations are:

- questionDelaySeconds - This sets how many seconds you want in between each question.
- imageRoundDelaySeconds - This sets how many seconds between the image round being sent out and regular trivia beginning.

