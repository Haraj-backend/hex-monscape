<script>
import { useRouter } from 'vue-router'

export default {
    setup() {
        const router = useRouter()

        return {
            router,
        }
    }
}
</script>

<template>
    <div id="welcome-screen-wrapper" class="flex relative w-screen sm:w-app h-app bg-cover">
        <div id="screen-wrapper">
            <!-- Game title -->
            <h1 class="game-title">
                HEX
                <br />MONSCAPE
            </h1>

            <!-- Initial actions -->
            <div class="game-initial-actions">
                <button
                    @click="router.push({ name: 'player-name' })"
                    class="bg-white rounded-lg shadow-[0_4px_0_rgba(0,0,0,0.2)] hover:shadow-[0_6px_0_rgba(0,0,0,0.2)] text-2xl max-[325px]:text-lg font-bold py-2 w-full sm:w-[264px] mt-10 sm:mt-0"
                >New Game</button>
                <button
                    @click="router.push({ name: 'about-screen' })"
                    class="bg-white rounded-lg shadow-[0_4px_0_rgba(0,0,0,0.2)] hover:shadow-[0_6px_0_rgba(0,0,0,0.2)] text-2xl max-[325px]:text-lg py-2 w-full sm:w-[264px]"
                >About</button>
            </div>
        </div>
    </div>
</template>