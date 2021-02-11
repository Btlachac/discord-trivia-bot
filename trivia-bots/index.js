require('dotenv').config();
const Discord = require("discord.js");
const botCommands = require('./commands.js').botCommands;


startTriviaBot(process.env.HOST_TOKEN);

console.log(botCommands);


async function startTriviaBot(token){

  const client = new Discord.Client();
  client.on('ready', () => {
    console.log(`Logged in as ${client.user.tag}!`);
  });
  
  
  client.on("message", async message => {
    // console.log(message);

    let cleanMessage = message.content.trim();

    //Finds the command where the command text matches the message
    const command = botCommands.find(c => c.commands.some(ct => ct.toUpperCase() === cleanMessage.toUpperCase()));
    if (command && command !== undefined){
      await command.function(client);
    } else {
      //invalid command
      //TODO: possibly send a message that the command isn't valid
    }


  });

  client.login(token);
}

// async function playTriviaRound(client, roundNumber, channelId){
//   const round = trivia.rounds.find(r => r.roundNumber === roundNumber);
//   const channel = client.channels.cache.get(channelId);

//   channel.send(`We will now begin **Round ${roundNumber}**`);
//   await sleep(1);
//   if(round.theme && round.theme.length > 0){
//     channel.send(`The theme for this round is **${round.theme}**`);
//     await sleep(1);
//   }
//   if (round.themeDescription && round.themeDescription.length > 0){
//     channel.send(round.themeDescription)
//   }
//   await sleep(15);

//   var i;
//   for(i=1; i< round.questions.length + 1; i++){
//     var q = round.questions[i-1];
//     channel.send(`**Question ${q.questionNumber}: ** ${q.question}`);
//     if (roundNumber != 5){
//       await sleep(config.questionDelaySeconds);

//     } else{
//       await sleep(5);
//     }
//   }
// }

// async function playTrivia(client, fromRoundNumber, channelId){
//   for (var i = fromRoundNumber; i <= 6; i++){
//     await playTriviaRound(client, i, channelId);
//     await sleep(config.questionDelaySeconds);
//   }
// }


// async function sleep(s){
//   return new Promise(resolve => setTimeout(resolve,s * 1000));
// }


// function startAudioBots(){
//   config.audioBotToChannelPairings.forEach(btcp => {
//     startAudioBot(btcp.token, btcp.channelId);
//   })

// }




// async function getNewTrivia(){
//   let response = await axios.get("http://localhost:8080/trivia");
//   return response.data;
//   // console.log(JSON.stringify(result.data, 0, 1));
// }
