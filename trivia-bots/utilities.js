
module.exports = {
    sleep: sleep,
    getTriviaChannel: getTriviaChannel,
    getBotCommandChannel: getBotCommandChannel
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