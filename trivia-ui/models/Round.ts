import Question from './Question';

export default class Round {
    roundNumber: Number;
    theme: string;
    questions: Question[];

    constructor(){
        this.roundNumber = 0;
        this.theme = '';
        this.questions = [];
    }

}