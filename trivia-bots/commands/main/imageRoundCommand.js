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
		});
	}

    async run(message) {
        // console.log('test');
        // return message.say('test');
        await botFunctions.sendImageRound(this.client);
        // const channel = utilities.getTriviaChannel(this.client);
        // channel.send(`**Image Round:**`)
        // channel.send(trivia.imageRoundTheme);
        // await utilities.sleep(1);
        // channel.send(trivia.imageRoundDetail);
        // await utilities.sleep(1);
        // channel.send(trivia.imageRoundURL); 
        return message.say('test')
    }
};