const { Command } = require('discord.js-commando');
const botFunctions = require('../../botFunctions')

module.exports = class MarkUsedCommand extends Command {
	constructor(client) {
		super(client, {
			name: 'mark_used',
			aliases: ['used'],
			group: 'main',
			memberName: 'mark used',
			description: 'Mark current trivia as used',
            guildOnly: true
		});
	}

    async run(message) {
        await botFunctions.markTriviaUsed(this.client);
        return message.say('trivia has been marked as used')
    }
};