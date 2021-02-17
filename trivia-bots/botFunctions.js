const fs = require('fs');
const utilities = require('./utilities')
const axios = require('axios');
const config = require('./config.json');

var trivia = null;

const AUDIO_FILE_LOCATION = "./audio.mp3"


module.exports = {
  startTrivia: startTrivia,
  stopTrivia: stopTrivia,
  pauseTrivia: pauseTrivia,
  startAudioRound: startAudioRound,
  sendImageRound: sendImageRound,
  sendAnswerSheet: sendAnswerSheet,
  getNextTrivia: getNextTrivia,
  markTriviaUsed: markTriviaUsed,
  resumeTrivia: resumeTrivia,
}

// TODO: any method using trivia needs to null check it and send a message if we don't have a valid trivia yet

//TODO: add another command/function for running just core trivia

//TODO: wrap up image, audio, and regular trivia all in this one cfunction
async function startTrivia(client) {
  //If we don't currently have a trivia, get a new one
  if (trivia === null){
    await getNextTrivia(client)
  }

  await sendImageRound(client);
  await utilities.sleep(config.imageRoundDelaySeconds);

  for (var i=1; i <= 6; i++){
    await playTriviaRound(client, i);
    await utilities.sleep(config.questionDelaySeconds);
  }

  await startAudioRound(client);

  // TODO: at the end of a trivia we should mark it as used in the DB and possibly null it
  //Reasoning is that the user only erver has to type start/start_trivia each week and everything else is hidden from them
}

async function playTriviaRound(client, roundNumber) {
  const round = trivia.rounds.find(r => r.roundNumber === roundNumber);
  const channel = utilities.getTriviaChannel(client);

  channel.send(`We will now begin **Round ${roundNumber}**`);
  await utilities.sleep(1);
  if (round.theme && round.theme.length > 0) {
    channel.send(`The theme for this round is **${round.theme}**`);
    await utilities.sleep(1);
  }
  if (round.themeDescription && round.themeDescription.length > 0) {
    channel.send(round.themeDescription)
  }
  await utilities.sleep(15);

  var i;
  for (i = 1; i < round.questions.length + 1; i++) {
    var q = round.questions[i - 1];
    channel.send(`**Question ${q.questionNumber}: ** ${q.question}`);
    if (roundNumber != 5) {
      await utilities.sleep(config.questionDelaySeconds);

    } else {
      await utilities.sleep(5);
    }
  }
}

async function stopTrivia(client) {

}

async function pauseTrivia(client) {

}

async function startAudioRound(client) {
  const channel = utilities.getTriviaChannel(client);
  channel.send(`It's now time for the Audio round.`);
  await utilities.sleep(1);
  channel.send(trivia.audioRoundTheme);
  await utilities.sleep(1);
  channel.send(`One of our lovely audio bots will join your channel in just a moment to play the clips, Good Luck!`);

  let audioBotToChannelPairings = unflattenAudioBotChannelPairings();

  writeAudioFile();

  await startAudioBots(audioBotToChannelPairings);
}

async function sendImageRound(client) {
  const channel = utilities.getTriviaChannel(client);
  channel.send(`**Image Round:**`)
  channel.send(trivia.imageRoundTheme);
  await utilities.sleep(1);
  channel.send(trivia.imageRoundDetail);
  await utilities.sleep(1);
  channel.send(trivia.imageRoundURL);
}

async function sendAnswerSheet(client) {
  const channel = utilities.getTriviaChannel(client);
  channel.send('**Answer Sheet:**')
  channel.send(trivia.answersURL);
}

async function getNextTrivia(client) {
  const channel = utilities.getBotCommandChannel(client);
  let baseUrl = process.env.API_URL;
  let response = await axios.get(`${baseUrl}/trivia`);
  trivia = response.data;
  console.log(JSON.stringify(trivia))

  channel.send("Successfully retrieved new trivia");

}

async function markTriviaUsed(client) {

}

async function resumeTrivia(client) {

}


async function startAudioBot(token, channelId) {
  const client = new Discord.Client();
  client.on('ready', async () => {
    const VC = client.channels.cache.get(channelId);
    await utilities.sleep(30);
    VC.join().then(async connection => {
      await utilities.sleep(2);
      await connection.play(config.audioFileName);

      //Bot leaves channel after 6 minutes
      await utilities.sleep(360);
      voiceChannel.leave();


    }).catch(err => {
      console.log(err);
    })
  });

  client.login(token);
}

function unflattenAudioBotChannelPairings() {
  let audioBotToChannelPairings = []

  let currentToken = process.env.AUDIO_BOT_TOKEN_1;
  let currentChannelId = process.env.AUDIO_BOT_CHANNEL_ID_1;

  let i = 1;

  while (currentToken && currentToken.length > 0 && currentChannelId && currentChannelId.length > 0) {
    let newPairing = {
      token: currentToken,
      channelId: currentChannelId
    }

    audioBotToChannelPairings.push(newPairing);
    i++;
    currentToken = process.env[`AUDIO_BOT_TOKEN_${i}`];
    currentChannelId = process.env[`AUDIO_BOT_CHANNEL_ID_${i}`];
  }

  return audioBotChannelPairings;
}

function writeAudioFile() {
  //Delete audio file if the old one is there
  fs.unlink(AUDIO_FILE_LOCATION);

  //Decode audio binary
  let decodedAudioBinary = atob(trivia.audioBinary);

  //write new audio file
  fs.writeFile(AUDIO_FILE_LOCATION, decodedAudioBinary);
}

async function startAudioBots(botChannelPairings) {
  for (botChannelPairing of botChannelPairings) {
    await startAudioBot(token, channelId)
  }
}


