<template>
  <v-form @submit.prevent="submit">
    <v-text-field v-model="answersURL" label="Answer Sheet URL"> </v-text-field>
    <v-text-field
      v-model="imageRoundURL"
      label="Image Round URL"
    ></v-text-field>
    <v-text-field
      v-model="audioRoundTheme"
      label="Audio Round Theme"
    ></v-text-field>

    <v-file-input
      v-model="triviaFile"
      label="Trivia Excel File"
      @change="triviaFileUploaded()"
    ></v-file-input>

    <v-file-input
      v-model="audioFile"
      label="Audio Round File"
      @change="audioFileUploaded()"
    ></v-file-input>

    <v-btn @click="submit()"> Submit </v-btn>
  </v-form>
</template>
<script>
const XLSX = require("xlsx");
export default {
  data: () => ({
    triviaFile: null,
    triviaData: {},
    imageRoundURL: null,
    answersURL: null,
    audioRoundTheme: null,
    audioFile: null,
    audioBinary: null,
  }),
  methods: {
    audioFileUploaded() {
      console.log('audio uploaded')
      var reader = new FileReader();
      reader.readAsBinaryString(this.audioFile);
      reader.onload = (e) => {
        console.log('audio loaded')
        //TODO: rename variables and make this process more clear

        this.audioBinary = btoa(e.target.result);
      };
    },
    //TODO: Trim whitespace on everything - should help prevent the startsWith from breaking
    triviaFileUploaded() {
      var reader = new FileReader();

      reader.readAsBinaryString(this.triviaFile);

      reader.onload = (e) => {
        console.log(e);
        var data = e.target.result;
        var workbook = XLSX.read(data, {
          type: "binary",
        });

        let rawTriviaJSON = XLSX.utils.sheet_to_json(workbook.Sheets["Q+A"], {
          header: "A",
        });
        this.formatTriviaData(rawTriviaJSON);
      };
    },
    formatTriviaData(rows) {

      rows = convertRowsToSringsAndTrim(rows)

      this.triviaData.rounds = [];

      var i = 0;
      var roundNumber = 0;
      var questionNumber = 0;


      for (i = 0; i < 52; i++) {
        var row = rows[i];
        if ((row.B && row.B.toUpperCase().includes("IMAGE ROUND")) || (row.A && row.A.toUpperCase() == "ROUND 4")) {
          this.triviaData.imageRoundTheme = trimRoundTheme(row.B);
          i++;
          this.triviaData.imageRoundDetail = rows[i].B;


          //We found the image round so advance counter until the start of the next round
          while (((!rows[i+1].A || !rows[i+1].A.toUpperCase().startsWith("ROUND")) && (!rows[i+1].B || !rows[i+1].B.toUpperCase().startsWith("ROUND"))) && i < rows.length){
            i++;
          }

        } else if ((row.B && row.B.toUpperCase().startsWith("ROUND") && !row.B.toUpperCase().includes("IMAGE ROUND")) || row.A && row.A.toUpperCase().startsWith("ROUND")) {
          questionNumber = 0;
          roundNumber++;

          this.triviaData.rounds[roundNumber - 1] = {
            roundNumber: roundNumber,
            theme: trimRoundTheme(row.B),
            questions: [],
          };
        } else if ((!row.A || row.A == "") && (!row.C || row.C == "")) {
          this.triviaData.rounds[roundNumber - 1].themeDescription = row.B;
        } else if (row.B && row.B != "") {
          questionNumber++;
          this.triviaData.rounds[roundNumber - 1].questions[
            questionNumber - 1
          ] = {
            questionNumber: questionNumber,
            question: row.B,
          };
        }
      }

      const jsonString = JSON.stringify(this.triviaData, null, 1);

      console.log(jsonString);

      // this.triviaData = jsonString;
    },
    submit() {
      this.triviaData.imageRoundURL = this.imageRoundURL;
      this.triviaData.audioRoundTheme = this.audioRoundTheme;
      this.triviaData.answersURL = this.answersURL;
      this.triviaData.audioBinary = this.audioBinary;
      console.log(JSON.stringify(this.triviaData));

      this.$axios.$post('trivia', this.triviaData);
    },

  },
};

function trimRoundTheme(roundTheme) {
  if (roundTheme) {
    roundTheme = roundTheme.split(/:(.+)/)[1];
  }
  return roundTheme;
}

function convertRowsToSringsAndTrim(rows){
  for (let row of rows){
    for (let prop in row){
      row[prop] = row[prop].toString().trim();
    }
  }
  return rows;
}
</script>
<style>
</style>
