const axios = require('axios');

module.exports = {
    getNextTrivia: getNextTrivia,
    markTriviaUsed: markTriviaUsed,
  }

async function getNextTrivia() {
    let baseUrl = process.env.API_URL;
    try {
        let response = await axios.get(`${baseUrl}/trivia`);
        return response.data;
    } catch (e){
        console.log(e);
        return null;
    }

}

async function markTriviaUsed(triviaId){
    let baseUrl = process.env.API_URL;
    try {
        await axios.put(`${baseUrl}/trivia/${triviaId}/mark-used`);
        return true;
    } catch(e){
        return false;
    }
}