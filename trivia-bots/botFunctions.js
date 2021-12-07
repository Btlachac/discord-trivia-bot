const fs = require('fs');
const utilities = require('./utilities')
const config = require('./config.json');
const Discord = require('discord.js')
const triviaService = require('./services/triviaService')

var trivia = null;

const AUDIO_FILE_LOCATION = "./audio.ogg"


module.exports = {
  startTrivia: startTrivia,
  startAudioRound: startAudioRound,
  sendImageRound: sendImageRound,
  sendAnswerSheet: sendAnswerSheet,
  getNextTrivia: getNextTrivia,
  markTriviaUsed: markTriviaUsed,
  restartBot: restartBot
}

// TODO: any method using trivia needs to null check it and send a message if we don't have a valid trivia yet

//TODO: add another command/function for running just core trivia

//TODO: wrap up image, audio, and regular trivia all in this one cfunction
async function startTrivia(client) {
  await getNextTrivia()

  await sendTriviaOverview(client);
  await utilities.sleep(20);

  await sendImageRound(client);
  await utilities.sleep(config.imageRoundDelaySeconds);

  for (var i=1; i <= trivia.rounds.length; i++){
    await playTriviaRound(client, i);
    await utilities.sleep(config.questionDelaySeconds);
  }

  sendMegaRoundReminder(client);
  await utilities.sleep(10);

  if (trivia.audioBinary){
    await startAudioRound(client);
  }

  await markTriviaUsed();
}

async function playTriviaRound(client, roundNumber) {
  const round = trivia.rounds.find(r => r.roundNumber === roundNumber);
  const channel = utilities.getTriviaChannel(client);

  let questionDelay = config.questionDelaySeconds;

  if (round.roundType.name.toUpperCase() == "LIGHTNING"){
    questionDelay = 5;
  }


  channel.send(`We will now begin **Round ${roundNumber}**`);
  await utilities.sleep(1);

  if (round.roundType.name.toUpperCase() == "LIST"){
    channel.send("This is a List Round! It's one question with multiple correct answers... teams get 1 point for each correct answer they write down. ");
    await utilities.sleep(1);
  }
  if (round.theme && round.theme.length > 0) {
    channel.send(`The theme for this round is **${round.theme}**`);
    await utilities.sleep(1);
  }
  if (round.themeDescription && round.themeDescription.length > 0) {
    channel.send(round.themeDescription)
  }

  await utilities.sleep(15);

  round.questions.sort((a,b) => (a.questionNumber > b.questionNumber) ? 1 : -1);

  if (round.roundType.name.toUpperCase() == "LIST"){
    var q = round.questions[0];
    channel.send(`**Question: ** ${q.question}`);
    await utilities.sleep(questionDelay * 2);
  }
  else {
    var i;
    for (i = 1; i < round.questions.length + 1; i++) {
      var q = round.questions[i - 1];
      channel.send(`**Question ${q.questionNumber}: ** ${q.question}`);
      await utilities.sleep(questionDelay);
    }
  }


}


async function startAudioRound(client) {
  // const channel = utilities.getTriviaChannel(client);
  // channel.send(`It's now time for the **Audio Round**`);
  // await utilities.sleep(1);
  // channel.send(trivia.audioRoundTheme);
  // await utilities.sleep(1);
  // channel.send(`One of our lovely audio bots will join your channel in just a moment to play the clips, Good Luck!`);

  let audioBotToChannelPairings = utilities.unflattenAudioBotChannelPairings();

  // writeAudioFile();

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

async function sendTriviaOverview(client) {
  const channel = utilities.getTriviaChannel(client);
  channel.send(`**Trivia Overview:**`)
  let overviewMessage = "Welcome to Trivia, tonight's game will include:\n"
  
  if (trivia.imageRoundURL && trivia.imageRoundURL.length > 0){
    overviewMessage += "- An Image Round\n"
  }

  overviewMessage += `- ${trivia.rounds.length} Rounds of regular trivia\n`

  if (trivia.audioBinary){
    overviewMessage += "- An Audio Round\n"
  }

  channel.send(overviewMessage);
}



function sendAnswerSheet(client) {
  const channel = utilities.getTriviaChannel(client);
  channel.send('**Answer Sheet:**')
  channel.send(trivia.answersURL);
}

function sendMegaRoundReminder(client) {
  const channel = utilities.getTriviaChannel(client);
  channel.send('**Reminder about the Mega Round:**')
  channel.send('You can pick any regular trivia round as your mega round, **excluding** the **list round**, the **image round** and the **audio round**. To do so you number your answers from 5-1 and you will get that many points if that answer is correct.')
}

async function getNextTrivia() {
  trivia = await triviaService.getNextTrivia();
}

async function startAudioBot(token, channelId) {
  const client = new Discord.Client();
  client.on('ready', async () => {
    const VC = client.channels.cache.get(channelId);
    // await utilities.sleep(30);
    VC.join().then(async connection => {
      // await utilities.sleep(2);
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

function restartBot(client) {
  client.destroy();
  client.login(process.env.HOST_TOKEN);
}


