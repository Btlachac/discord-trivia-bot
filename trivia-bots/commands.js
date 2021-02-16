const botFunctions = require('./botFunctions')
const utilities = require('./utilities')


const botCommands = [
    {
        commands: ["help"],
        helpText: "Displays all available commands",
        function: sendHelpMessage
    },
    {
        commands: ["start", "start_trivia"],
        helpText: "Starts trivia, sending out the image round, then 6 rounds of trivia, then the audio round",
        function: botFunctions.startTrivia
    },
    // {
    //     commands: ["stop"],
    //     helpText: "Stops the current trivia",
    //     function: botFunctions.stopTrivia
    // },
    // {
    //     commands: ["pause"],
    //     helpText: "Pauses the current trivia",
    //     function: botFunctions.pauseTrivia
    // },
    // {
    //     commands: ["resume"],
    //     helpText: "Resume trivia from pause",
    //     function: botFunctions.resumeTrivia
    // },
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
    //TODO: need status command to check if we currently have a trivia etc

]


// TODO: need to do something about this inconsistency
//Either all these functions need their base call in here or something else

async function sendHelpMessage(client){
    const channel = utilities.getBotCommandChannel(client);
    
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

    

