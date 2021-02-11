const fs = require('fs');
const utilities = require('./utilities')

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


async function startTrivia(client) {

}

async function stopTrivia(client) {

}

async function pauseTrivia(client) {

}

async function startAudioRound(client) {
    const channel = getTriviaChannel(client);
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
    const channel = getTriviaChannel(client);
    channel.send(trivia.imageRoundTheme);
    await utilities.utilities.sleep(1);
    channel.send(trivia.imageRoundDetail);
    await utilities.utilities.sleep(1);
    channel.send(trivia.imageRoundURL);
}

async function sendAnswerSheet(client) {
    const channel = getTriviaChannel(client);
    channel.send('**Answer Sheet:**')
    channel.send(trivia.answersURL);
}

async function getNextTrivia(client) {
    const channel = getTriviaChannel(client);
    let baseUrl = process.env.API_URL;
    let response = await axios.get(`${baseUrl}/trivia`);
    trivia = response.data;

    channel.send("Successfully retrieved now trivia")
    
}

async function markTriviaUsed(client) {

}

async function resumeTrivia(client) {

}




async function startAudioBot(token, channelId){
    const client = new Discord.Client();
    client.on('ready', async () => {
        const VC = client.channels.cache.get(channelId);
        await utilities.sleep(30);
        VC.join().then(async connection =>
        {
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
    let audioBotChannelPairings = []

    let currentToken = process.env.AUDIO_BOT_TOKEN_1;
    let currentChannelId = process.env.AUDIO_BOT_CHANNEL_ID_1;

    let i = 1;

    while(currentToken && currentToken.length > 0 && currentChannelId && currentChannelId.length > 0){
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

    fs.writeFile(AUDIO_FILE_LOCATION, decodedAudioBinary);
}

async function startAudioBots(botChannelPairings){
    for (botChannelPairing of botChannelPairings) {
        await startAudioBot(token, channelId)
    }
}

function getTriviaChannel(client){
    return client.channels.cache.get(process.env.TRVIA_CHANNEL_ID);
}
