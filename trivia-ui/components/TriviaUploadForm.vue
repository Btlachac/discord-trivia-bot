<template>
  <v-form @submit.prevent="submit">
    <v-text-field v-model="answersURL" label="Answer Sheet URL"> </v-text-field>
    <v-text-field
      v-model="imageRoundURL"
      label="Image Round URL"
    ></v-text-field>
    <!-- <v-text-field
      v-model="audioRoundTheme"
      label="Audio Round Theme"
    ></v-text-field> -->

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
    <!-- <div>
      <pre>{{ triviaString }}</pre>
    </div> -->

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
    // triviaString: null
  }),
  methods: {
    audioFileUploaded() {
      console.log("audio uploaded");
      var reader = new FileReader();
      reader.readAsBinaryString(this.audioFile);
      reader.onload = (e) => {
        console.log("audio loaded");
        //TODO: rename variables and make this process more clear

        this.audioBinary = btoa(e.target.result);
      };
    },
    triviaFileUploaded() {
      var reader = new FileReader();

      reader.readAsBinaryString(this.triviaFile);

      reader.onload = (e) => {
        console.log(e);
        var data = e.target.result;
        var workbook = XLSX.read(data, {
          type: "binary",
        });
        console.log(workbook.SheetNames);
        let sourceSheet = workbook.SheetNames.find(
          (sn) => sn.toUpperCase().trim() === "SOURCE"
        );

        if (sourceSheet) {
          let rawTriviaJSON = XLSX.utils.sheet_to_json(workbook.Sheets[sourceSheet]);
          this.parseTriviaFromSource(rawTriviaJSON);
        } else {
          let rawTriviaJSON = XLSX.utils.sheet_to_json(workbook.Sheets["Q+A"], {
            header: "A",
          });
          this.parseTriviaFromQA(rawTriviaJSON);
        }
      };
    },
    parseTriviaFromQA(rows) {
      rows = convertRowsToSringsAndTrim(rows);

      this.triviaData.rounds = [];

      var i = 0;
      var roundNumber = 0;
      var questionNumber = 0;

      if (rows[53] && rows[53].B) {
        this.audioRoundTheme = rows[53].B;
      }

      for (i = 0; i < 52; i++) {
        var row = rows[i];
        if (
          (row.B && row.B.toUpperCase().includes("IMAGE ROUND")) ||
          (row.A && row.A.toUpperCase() == "ROUND 4")
        ) {
          this.triviaData.imageRoundTheme = trimRoundTheme(row.B);
          i++;
          this.triviaData.imageRoundDetail = rows[i].B;

          //We found the image round so advance counter until the start of the next round
          while (
            (!rows[i + 1].A ||
              !rows[i + 1].A.toUpperCase().startsWith("ROUND")) &&
            (!rows[i + 1].B ||
              !rows[i + 1].B.toUpperCase().startsWith("ROUND")) &&
            i < rows.length
          ) {
            i++;
          }
        } else if (
          (row.B &&
            row.B.toUpperCase().startsWith("ROUND") &&
            !row.B.toUpperCase().includes("IMAGE ROUND")) ||
          (row.A && row.A.toUpperCase().startsWith("ROUND"))
        ) {
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

      // const jsonString = JSON.stringify(this.triviaData, null, 1);

      // console.log(jsonString);

      // this.triviaString = jsonString;
    },
    parseTriviaFromSource(rows) {
      console.log('Using source sheet')

      rows = convertRowsToSringsAndTrim(rows);

      // console.log(JSON.stringify(rows, 0 ,1))

      this.triviaData.rounds = [];

      //filter out empty rows
      let rowsWithData = rows.filter(r => (r.question && r.question.length > 0) || (r.answer && r.answer.length > 0));

      //find a row from the image round to get the image round title and theme
      let imageRoundRow = rowsWithData.find(r => r.round_type.toUpperCase() === "IMAGE ROUND");
      
      if (imageRoundRow && imageRoundRow !== undefined){
          this.triviaData.imageRoundTheme = imageRoundRow.round_title;
          this.triviaData.imageRoundDetail = imageRoundRow.round_description;
      }

      //find a row from the audio round to get the audio round theme
      let audioRoundRow = rowsWithData.find(r => r.round_type.toUpperCase() === "SOUND ROUND");

      if (audioRoundRow && audioRoundRow !== undefined){
          this.triviaData.audioRoundTheme = audioRoundRow.round_description;
      }

      var roundNumber = 0;
      var questionNumber = 0;

      //filter out all non-regular trivia rounds
      let regularTriviaRows = rowsWithData.filter(r => r.round_type.toUpperCase() !== "IMAGE ROUND" && r.round_type.toUpperCase() !== "TIEBREAKER" && r.round_type.toUpperCase() !== "SOUND ROUND")

      //loop through trivia rounds and transform data into our json object
      for (const row of regularTriviaRows) {
        //this is the start of a new round
        if (row.question_number === "1") {
          roundNumber++;
          questionNumber = 1;

          let newRound = {
            roundNumber: roundNumber,
            theme: row.round_title,
            questions: [],
          };

          if (row.round_description && row.round_description.length > 0) {
            newRound.themeDescription = row.round_description;
          }

          this.triviaData.rounds.push(newRound);
        }

        let newQuestion = {
          questionNumber: questionNumber,
          question: row.question,
        };

        this.triviaData.rounds[roundNumber - 1].questions.push(newQuestion);
        questionNumber++;

      }

      // console.log(JSON.stringify(this.triviaData, 0, 1))

    },
    submit() {
      this.triviaData.imageRoundURL = this.imageRoundURL;
      this.triviaData.answersURL = this.answersURL;
      if (this.triviaData.audioBinary) {
        this.triviaData.audioBinary = this.audioBinary;
        this.triviaData.audioRoundTheme = this.audioRoundTheme;
      }
      console.log(JSON.stringify(this.triviaData));

      this.$axios.$post("trivia", this.triviaData);
    },
  },
};

function trimRoundTheme(roundTheme) {
  if (roundTheme) {
    roundTheme = roundTheme.split(/:(.+)/)[1];
  }
  return roundTheme;
}

function convertRowsToSringsAndTrim(rows) {
  for (let row of rows) {
    for (let prop in row) {
      row[prop] = row[prop].toString().trim();
    }
  }
  return rows;
}
</script>
<style>
</style>
