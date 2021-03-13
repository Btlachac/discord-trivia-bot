const { Command } = require('discord.js-commando');
const botFunctions = require('../../botFunctions')

module.exports = class StartTriviaCommand extends Command {
	constructor(client) {
		super(client, {
			name: 'answer_sheet',
			aliases: ['answers'],
			group: 'main',
			memberName: 'answer sheet',
			description: 'Send out a link to the answer sheet',
            guildOnly: true
		});
	}

    run(message) {
        botFunctions.sendAnswerSheet(this.client);
        return message.say('Answer Sheet Sent')
    }
};