<script setup lang="ts">
const playerName = ref("")
</script>
<template>
    <div id="welcome-screen-wrapper" class="flex relative w-app h-app bg-cover">
        <div class="p-[160px]">
            <!-- Game title -->
            <h1 class="game-title">NEW GAME</h1>

            <!-- Initial actions -->
            <div class="game-initial-actions">
                <label class="flex flex-col text-2xl" for="input_player-name">
                    <span class="mb-2">Who are you?</span>
                    <input v-model="playerName" type="text" id="input_player-name" />
                </label>
                <div class="flex">
                <NuxtLink to="/" class="button w-full">Choose your partner</NuxtLink>
            </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
#input_player-name {
  @apply outline outline-2 outline-[rgba(0,0,0,.1)] focus:outline-[rgba(0,0,0,.3)] rounded-lg;
  @apply py-2 px-4;
}
</style>