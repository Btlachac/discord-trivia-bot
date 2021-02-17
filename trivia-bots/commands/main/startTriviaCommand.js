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
		});
	}

    async run(message) {
        // console.log('test');
        // return message.say('test');
        await botFunctions.startTrivia(this.client);
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