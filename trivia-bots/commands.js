const botFunctions = require('./botFunctions')





const botCommands = [
    {
        commands: ["start", "start_trivia"],
        helpText: "Starts trivia, sending out the image round, then 6 rounds of trivia, then the audio round",
        function: botFunctions.startTrivia
    },
    {
        commands: ["help"],
        helpText: "Displays all available commands",
        function: sendHelpMessage
    },
    {
        commands: ["stop"],
        helpText: "Stops the current trivia",
        function: botFunctions.stopTrivia
    },
    {
        commands: ["pause"],
        helpText: "Pauses the current trivia",
        function: botFunctions.pauseTrivia
    },
    {
        commands: ["resume"],
        helpText: "Resume trivia from pause",
        function: botFunctions.resumeTrivia
    },
    {
        commands: ["image_round"],
        helpText: "Sends out the image round",
        function: botFunctions.sendImageRound
    },
    {
        commands: ["audio_round"],
        helpText: "Begins the audio round",
        function: botFunctions.startAudioRound
    },
    {
        commands: ["answers", "answer_sheet"],
        helpText: "Sends out a link to the answer sheet",
        function: botFunctions.sendAnswerSheet
    },
    {
        commands: ["next", "next_trivia"],
        helpText: "Fetches a new trivia",
        function: botFunctions.getNextTrivia
    },
    {
        commands: ["mark_used"],
        helpText: "Marks current trivia as used",
        function: botFunctions.markTriviaUsed
    },
    //TODO: make generic command for round_x
]

// TODO: need to do something about this inconsistency
//Either all these functions need their base call in here or something else

async function sendHelpMessage(client){
    const channel = client.channels.cache.get(process.env.BOT_COMMAND_CHANNEL_ID);

    let helpMessage = 'These are the available commands: \n ```'

    let commands = botCommands.map(bc => bc.commands.join(' or ') + ' - ' + bc.helpText)

    helpMessage += commands.join('\n');
    // const reducer = (accumulator, currentValue) => accumulator + currentValue + '\n';

    // helpMessage += + commands.reduce(reducer);

    helpMessage += '```';
    channel.send(helpMessage);
}



module.exports = {
    botCommands: botCommands
}

    

