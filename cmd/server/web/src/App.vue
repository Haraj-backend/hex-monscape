<script>
import { useStore } from "./store";
import { randomPick } from "./composables/utils";
import { onMounted } from "vue";
import { useRoute } from "vue-router";

export default {
  setup() {
    const store = useStore()
    const backgrounds = [
      'https://idev-images-test.s3.eu-west-1.amazonaws.com/hex-pokebattle-backgrounds/bulbasaur.jpg',
      'https://idev-images-test.s3.eu-west-1.amazonaws.com/hex-pokebattle-backgrounds/charmander.jpg',
      'https://idev-images-test.s3.eu-west-1.amazonaws.com/hex-pokebattle-backgrounds/charmeleon.jpg',
      'https://idev-images-test.s3.eu-west-1.amazonaws.com/hex-pokebattle-backgrounds/ivysaur.jpg',
      'https://idev-images-test.s3.eu-west-1.amazonaws.com/hex-pokebattle-backgrounds/pikachu.jpg',
      'https://idev-images-test.s3.eu-west-1.amazonaws.com/hex-pokebattle-backgrounds/squirtle.jpg'
    ]

    let selectedBg = store.gameBackground
    if (!selectedBg || selectedBg === '') {
      const d = new Date().getDate()
      selectedBg = randomPick(backgrounds)
      store.setGameBackground(selectedBg, d)
    }

    onMounted(() => {
      const gw = document.getElementById("game-wrapper")
      gw.style.backgroundImage = `url(${store.gameBackground})`
    })
  }
}
</script>

<template>
  <div id="game-wrapper">
    <router-view></router-view>
  </div>
</template>
