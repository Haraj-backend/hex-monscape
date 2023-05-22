<script setup lang="ts">
const utils = useUtils();

const backgrounds = [
  'https://idev-images-test.s3.eu-west-1.amazonaws.com/hex-pokebattle-backgrounds/bulbasaur.jpg',
  'https://idev-images-test.s3.eu-west-1.amazonaws.com/hex-pokebattle-backgrounds/charmander.jpg',
  'https://idev-images-test.s3.eu-west-1.amazonaws.com/hex-pokebattle-backgrounds/charmeleon.jpg',
  'https://idev-images-test.s3.eu-west-1.amazonaws.com/hex-pokebattle-backgrounds/ivysaur.jpg',
  'https://idev-images-test.s3.eu-west-1.amazonaws.com/hex-pokebattle-backgrounds/pikachu.jpg',
  'https://idev-images-test.s3.eu-west-1.amazonaws.com/hex-pokebattle-backgrounds/squirtle.jpg'
]

let selectedBg = utils.randomPick(backgrounds);

onMounted(() => {
  const gw = document.getElementById("game-wrapper");
  if (gw) {
    gw.style.backgroundImage = `url(${selectedBg})`
  }
})

</script>

<template>
  <div id="game-wrapper">
    <h1>Hex Pokebattle Nuxt</h1>
  </div>
</template>
