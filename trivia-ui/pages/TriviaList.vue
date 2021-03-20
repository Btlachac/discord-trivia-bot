<template>
  <div>
    <v-data-table
      :headers="headers"
      :items="triviaData"
      item-key="id"
      class="elevation-1"
    >
    </v-data-table>
  </div>
</template>

<script>

export default {
 data () {
      return {
        headers: [
          { text: 'Date Created', value: 'dateCreated' },
          { text: 'Used', value: 'used' },
          { text: 'Date Used', value: 'dateUsed' },
          { text: 'Image Round Theme', value: 'imageRoundTheme' },
        ],
      }
    },
}
</script>
