const { Command } = require('discord.js-commando');
const botFunctions = require('../../botFunctions')

module.exports = class ImageRoundCommand extends Command {
	constructor(client) {
		super(client, {
			name: 'next_trivia',
			aliases: ['next'],
			group: 'main',
			memberName: 'next trivia',
			description: 'Fetches a new trivia',
            guildOnly: true
		});
	}

    async run(message) {
        await botFunctions.getNextTrivia();
        return message.say('New Trivia Fetched')
    }
};