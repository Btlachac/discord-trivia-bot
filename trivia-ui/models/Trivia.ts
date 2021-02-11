import Round from './Round';

export default class Trivia {
    rounds: Round[];
    imageRoundTheme: String;
    imageRoundDetail: String;
    imageRoundURL: String;
    answersURL: String;
    audioRoundTheme: String;
    audioBinary: String;

    constructor(){
        this.rounds = [];
        this.imageRoundTheme = '';
        this.imageRoundDetail = '';
        this.imageRoundURL = '';
        this.answersURL = '';
        this.audioRoundTheme = '';
        this.audioBinary = '';
    }

}