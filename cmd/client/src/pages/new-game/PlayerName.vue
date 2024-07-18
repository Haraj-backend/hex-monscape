<script>
import { ref } from "vue"
import { useRouter } from "vue-router"
import { useStore } from "../../store"

export default {
    setup() {
        const router = useRouter()
        const store = useStore()
        const playerName = ref("")

        if (store.getPlayerName) {
            playerName.value = store.getPlayerName
        }

        const choosePartner = () => {
            if (!playerName.value) return

            // update `playerName` in store
            store.setPlayerName(playerName.value)

            // route to the next page
            router.push({ name: 'choose-partner' })
        }

        return {
            playerName,
            choosePartner
        }
    }
}
</script>

<template>
    <div id="welcome-screen-wrapper" class="flex relative w-app h-app bg-cover">
        <div id="screen-wrapper">
            <!-- Game title -->
            <h1 class="game-title">NEW GAME</h1>

            <!-- Initial actions -->
            <div class="game-initial-actions">
                <label class="flex flex-col text-2xl" for="input_player-name">
                    <span class="mb-2">Who are you?</span>
                    <input v-model="playerName" type="text" id="input_player-name" class="max-w-[980px]" />
                </label>
                <button
                    @click="choosePartner"
                    class="bg-[rgba(0,0,0,.2)] rounded-lg text-2xl py-2 px-4 max-w-[980px]"
                >Choose your partner</button>
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