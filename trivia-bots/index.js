require('dotenv').config();
const { CommandoClient } = require('discord.js-commando');
const path = require('path');


const client = new CommandoClient({
  commandPrefix: '!',
  owner: process.env.USER_ID
});


main(process.env.HOST_TOKEN);

async function main(token){
  client.registry
	.registerDefaultTypes()
	.registerGroups([
		['main', 'main command group']
	])
	.registerDefaultGroups()
	.registerDefaultCommands()
	.registerCommandsIn(path.join(__dirname, 'commands'));


  client.once('ready', () => {
    console.log(`Logged in as ${client.user.tag}! (${client.user.id})`);
    client.user.setActivity('Trivia');
  });
  
  client.on('error', console.error);

  client.login(token);
}