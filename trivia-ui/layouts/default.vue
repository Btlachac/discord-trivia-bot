<template>
  <div class="outer-container">
    <div class="inner-container">
        <nav-bar class="nav-bar" />
        <Nuxt />
    </div>
  </div>
</template>

<script>

import NavBar from '~/components/NavBar.vue'

export default {
   components: {
    NavBar
  }
};

</script>

<style scoped>
.inner-container {
  max-width: 1100px;
  margin: 0 auto;
}

.nav-bar {
  margin-top: 20px;
  margin-left: auto;
  margin-right: auto;
}
</style>