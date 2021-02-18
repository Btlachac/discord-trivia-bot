const { Command } = require('discord.js-commando');
const botFunctions = require('../../botFunctions')

module.exports = class StartTriviaCommand extends Command {
	constructor(client) {
		super(client, {
			name: 'audio_round',
			aliases: ['audio'],
			group: 'main',
			memberName: 'audio round',
			description: 'Begins the audio round',
            guildOnly: true
		});
	}

    async run(message) {
        botFunctions.startAudioRound(this.client);
        return message.say('Audio round will begin shortly')
    }
};