<script>
import { useStore } from "../../store"
import PartnerCard from "../../components/PartnerCard.vue"
import { useRoute, useRouter } from "vue-router"
import { computed, onMounted, ref } from "vue"
import PokebattleHTTP from "../../composables/http_client"
import { turnStates } from "../../entity/game"

export default {
    components: {
        PartnerCard
    },
    setup() {
        // dependencies initialization
        const route = useRoute()
        const router = useRouter()
        const store = useStore()
        const client = new PokebattleHTTP()

        // reactive variables
        const currentGameData = computed(() => store.getGameData)
        const currentBattle = computed(() => currentGameData.value.scenario.split('_')[1])
        const battleState = computed(() => store.getBattleState)
        const gameFinished = computed(() => route.params.state === 'finished')
        const retry = ref(0)

        const getGameDetails = async () => {
            const resp = await client.getGameDetails(currentGameData.value.id)
            if (resp.ok) {
                store.setTheGame(resp.data)
            }
        }

        // methods
        const proceed = async (num) => {
            if (gameFinished.value) {
                router.push({ name: 'welcome-screen' })
            } else {
                // start the battle first and then redirect
                // to battle scene
                const battleResp = await client.startBattle(currentGameData.value.id)

                if (battleResp.ok) {
                    store.setTheBattle(battleResp.data)
                    router.push({ name: 'battle', params: { scenario: `scenario-${num}` } })
                } else {
                    if (retry.value < 3) {
                        // just create new game if it's error
                        // this is the BUG
                        // need to revisit the business logic
                        const newGameResp = await client.newGame({
                            PlayerName: store.getPlayerName,
                            PartnerID: store.getPartnerData.id
                        })

                        store.setTheGame(newGameResp.data)
                        retry.value++
                        // just try again start the battle
                        await proceed(num)
                    }
                }
            }
        }

        onMounted(() => {
            // update the game data if win
            // request game details
            if (battleState.value && battleState.value.state === turnStates.WIN) {
                getGameDetails()
            }
        })

        return {
            partner: store.partnerData,
            playerName: store.getPlayerName,
            gameFinished,
            currentGameData,
            currentBattle,
            proceed
        }
    }
}
</script>
<template>
    <div class="flex justify-between w-full h-app p-[160px]">
        <div class="left-side">
            <!-- Texts -->
            <p class="game-description text-3xl">
                {{ gameFinished ? 'Congratulations' : 'Hello' }},
                <span
                    class="player-name"
                >{{ playerName }}</span>
            </p>
            <p
                class="game-description text-xl mt-2"
            >{{ gameFinished ? 'You finished the game' : `Total wins ${currentGameData.battleWon}` }},</p>

            <!-- Battle scenario buttons -->
            <div class="battle-scenario-list flex flex-col gap-y-4 mt-24 pr-24">
                <p
                    class="game-description text-xl"
                >Choose the available scenario to start the battle and defeat the enemy!</p>
                <button
                    v-for="battleScenario in 3"
                    @click="proceed(battleScenario)"
                    class="w-[300px] rounded-lg text-2xl py-2 px-3"
                    :id="`battle-scenario-button-${battleScenario}`"
                    :class="battleScenario <= currentBattle ? 'battle-scenario-available' : 'battle-scenario-disabled'"
                >Scenario {{ battleScenario }}</button>
            </div>
        </div>
        <div class="right-side">
            <PartnerCard :partner="partner" />
        </div>
    </div>
</template>

<style scoped>
.battle-scenario-available {
    @apply bg-[rgba(0,0,0,.5)] hover:shadow-[0_4px_0_rgba(0,0,0,.65)] active:shadow-[0_-4px_0_rgba(0,0,0,.65)];
    @apply font-bold text-white;
}

.battle-scenario-disabled {
    @apply bg-[rgba(0,0,0,.25)] text-[rgba(0,0,0,.35)] cursor-not-allowed;
}
</style>