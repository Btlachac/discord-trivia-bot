
module.exports = {
    sleep: sleep,
    getTriviaChannel: getTriviaChannel,
    getBotCommandChannel: getBotCommandChannel,
    unflattenAudioBotChannelPairings: unflattenAudioBotChannelPairings
}

async function sleep(s){
    return new Promise(resolve => setTimeout(resolve,s * 1000));
}

function getTriviaChannel(client){
    return client.channels.cache.get(process.env.TRVIA_CHANNEL_ID);
}

function getBotCommandChannel(client){
    return client.channels.cache.get(process.env.BOT_COMMAND_CHANNEL_ID);
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
  
    return audioBotToChannelPairings;
  }