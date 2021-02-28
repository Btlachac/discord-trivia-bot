const fs = require('fs');
const utilities = require('./utilities')
const config = require('./config.json');
const Discord = require('discord.js')
const triviaService = require('./services/triviaService')

var trivia = null;

const AUDIO_FILE_LOCATION = "./audio.mp3"


module.exports = {
  startTrivia: startTrivia,
  startAudioRound: startAudioRound,
  sendImageRound: sendImageRound,
  sendAnswerSheet: sendAnswerSheet,
  getNextTrivia: getNextTrivia,
  markTriviaUsed: markTriviaUsed,
}

// TODO: any method using trivia needs to null check it and send a message if we don't have a valid trivia yet

//TODO: add another command/function for running just core trivia

//TODO: wrap up image, audio, and regular trivia all in this one cfunction
async function startTrivia(client) {
  //If we don't currently have a trivia, get a new one
  if (trivia === null){
    await getNextTrivia()
  }

  await sendImageRound(client);
  await utilities.sleep(config.imageRoundDelaySeconds);

  for (var i=1; i <= 6; i++){
    await playTriviaRound(client, i);
    await utilities.sleep(config.questionDelaySeconds);
  }

  await startAudioRound(client);

  await markTriviaUsed();
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

async function startAudioRound(client) {
  const channel = utilities.getTriviaChannel(client);
  channel.send(`It's now time for the **Audio Round**`);
  await utilities.sleep(1);
  channel.send(trivia.audioRoundTheme);
  await utilities.sleep(1);
  channel.send(`One of our lovely audio bots will join your channel in just a moment to play the clips, Good Luck!`);

  let audioBotToChannelPairings = utilities.unflattenAudioBotChannelPairings();

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

async function getNextTrivia() {
  trivia = await triviaService.getNextTrivia();
}

async function startAudioBot(token, channelId) {
  const client = new Discord.Client();
  client.on('ready', async () => {
    const VC = client.channels.cache.get(channelId);
    await utilities.sleep(30);
    VC.join().then(async connection => {
      await utilities.sleep(2);
      await connection.play(AUDIO_FILE_LOCATION);

      //Bot leaves channel after 6 minutes
      await utilities.sleep(360);
      VC.leave();
      client.destroy();

    }).catch(err => {
      console.log(err);
    })
  });

  client.login(token);
}



function writeAudioFile() {
  //Decode audio binary
  let decodedAudioBinary = Buffer.from(trivia.audioBinary, 'base64');

  //write new audio file - if a file exists it will be overwritten
  fs.writeFileSync(AUDIO_FILE_LOCATION, decodedAudioBinary);
}

async function startAudioBots(botChannelPairings) {
  for (botChannelPairing of botChannelPairings) {
    await startAudioBot(botChannelPairing.token, botChannelPairing.channelId)
  }
}

async function markTriviaUsed(){
  if (trivia){
    let response = await triviaService.markTriviaUsed(trivia.id);
  }
}


