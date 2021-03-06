const { Command } = require('discord.js-commando');
const botFunctions = require('../../botFunctions')

module.exports = class ImageRoundCommand extends Command {
	constructor(client) {
		super(client, {
			name: 'image_round',
			aliases: ['image'],
			group: 'main',
			memberName: 'image round',
			description: 'Sends out the image round',
            guildOnly: true
		});
	}

    async run(message) {
        botFunctions.sendImageRound(this.client);
        return message.say('Image Round Sent')
    }
};