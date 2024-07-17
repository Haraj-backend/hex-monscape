<script>
import { useStore } from "./store";
import { randomPick } from "./composables/utils";
import { onMounted } from "vue";

export default {
  setup() {
    const store = useStore()
    const backgrounds = [
      'https://haraj-sol-dev.s3.eu-west-1.amazonaws.com/hex-monscape/backgrounds/battle_base.jpg'
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