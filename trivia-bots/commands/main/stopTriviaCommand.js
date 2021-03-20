const { Command } = require('discord.js-commando');
const botFunctions = require('../../botFunctions')

module.exports = class StartTriviaCommand extends Command {
	constructor(client) {
		super(client, {
			name: 'stop_trivia',
			aliases: ['stop'],
			group: 'main',
			memberName: 'stops trivia',
			description: 'Stops trivia by restarting the bot',
            guildOnly: true
		});
	}

    run(message) {
        botFunctions.restartBot(this.client);
        return message.say('Trivia Stopped')
    }
};