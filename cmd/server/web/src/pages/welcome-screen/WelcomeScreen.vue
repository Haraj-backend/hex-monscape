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
    <div id="welcome-screen-wrapper" class="flex relative w-app h-app bg-cover">
        <div class="p-[160px]">
            <!-- Game title -->
            <h1 class="game-title">
                HEX
                <br />POKEBATTLE
            </h1>

            <!-- Initial actions -->
            <div class="game-initial-actions">
                <button
                    @click="router.push({ name: 'player-name' })"
                    class="bg-white rounded-lg shadow-[0_4px_0_rgba(0,0,0,0.2)] hover:shadow-[0_6px_0_rgba(0,0,0,0.2)] text-2xl font-bold py-2 w-[264px]"
                >New Game</button>
                <button
                    @click="router.push({ name: 'about-screen' })"
                    class="bg-white rounded-lg shadow-[0_4px_0_rgba(0,0,0,0.2)] hover:shadow-[0_6px_0_rgba(0,0,0,0.2)] text-2xl py-2 w-[264px]"
                >About</button>
            </div>
        </div>
    </div>
</template>