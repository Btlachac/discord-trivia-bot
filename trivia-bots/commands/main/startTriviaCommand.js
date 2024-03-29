const { Command } = require('discord.js-commando');
const botFunctions = require('../../botFunctions')

module.exports = class StartTriviaCommand extends Command {
	constructor(client) {
		super(client, {
			name: 'start_trivia',
			aliases: ['start'],
			group: 'main',
			memberName: 'start trivia',
			description: 'Starts trivia by sending out the image round, main trivia, and the audio round',
            guildOnly: true
		});
	}

    async run(message) {
        botFunctions.startTrivia(this.client);
        return message.say('Trivia Started')
    }
};