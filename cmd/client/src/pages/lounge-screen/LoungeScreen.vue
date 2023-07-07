<script>
import { useRouter } from "vue-router"
import { computed, onMounted } from "vue"
import { useStore } from "../../store"
import { getGameScenario, turnStates } from "../../entity/game"
import PartnerCard from "../../components/PartnerCard.vue"
import MonscapeHTTP from "../../composables/http_client"

export default {
    components: {
        PartnerCard
    },
    setup() {
        // dependencies initialization
        const router = useRouter()
        const store = useStore()
        const client = new MonscapeHTTP()

        // reactive variables
        const currentGameData = computed(() => store.getGameData)
        const gameFinished = computed(() => currentGameData.value.scenario === turnStates.END_GAME)
        const loungeOrder = computed(() => {
            switch (currentGameData.value.scenario) {
                case 'BATTLE_2':
                    return 'Win your 2nd battle to progress the game!';
                case 'BATTLE_3':
                    return 'Win your 3rd battle to progress the game!';
                case turnStates.END_GAME:
                    return 'You may continue doing battle or start new game!';
                default:
                    return 'Win your 1st battle to progress the game!';
            }
        })

        // methods
        const proceed = async () => {
            // start the battle first and
            // then redirect to battle scene
            const battleResp = await client.startBattle(currentGameData.value.id)
            if (battleResp.ok) {
                store.setTheBattle(battleResp.data)
                router.push({ name: 'battle', params: { scenario: `scenario-${getGameScenario(currentGameData.value)}` } })
            } else {
                console.error({ battleResp })
            }
        }

        const newGame = () => {
            // reset the game
            store.resetGame()
            router.push({ name: 'welcome-screen' })
        }

        const getGameDetails = async () => {
            const resp = await client.getGameDetails(currentGameData.value.id)
            if (resp.ok) {
                store.setTheGame(resp.data)
            } else {
                // reset the game and battle data from client storage
                // redirect to welcome screen
                newGame()
            }
        }

        onMounted(() => {
            // update the game data if win
            // request game details
            getGameDetails()
        })

        return {
            partner: store.partnerData,
            playerName: store.getPlayerName,
            gameFinished,
            currentGameData,
            loungeOrder,
            newGame,
            proceed
        }
    }
}
</script>
<template>
    <div class="flex justify-between w-full h-app p-[160px]">
        <div class="left-side pr-24">
            <!-- Texts -->
            <p class="game-description text-3xl">
                {{ gameFinished ? 'Congratulations' : 'Hello' }},
                <span
                    class="player-name"
                >{{ playerName }}</span>
            </p>
            <p
                class="game-description text-xl mt-2"
            >{{ gameFinished ? 'You have won the game!' : `Total wins ${currentGameData.battleWon}` }}</p>

            <!-- Battle scenario buttons -->
            <div class="battle-scenario-list flex flex-col gap-y-4 mt-24">
                <p class="game-description text-xl">{{ loungeOrder }}</p>
                <button
                    @click="proceed(battleScenario)"
                    class="bg-red-600 text-white w-[300px] rounded-lg text-2xl py-2 px-3"
                >Go to battle</button>
                <button
                    v-if="gameFinished"
                    @click="newGame"
                    class="bg-[rgba(0,0,0,.25)] text-white w-[300px] rounded-lg text-2xl py-2 px-3"
                >New Game</button>
            </div>
        </div>
        <div class="right-side">
            <PartnerCard :partner="partner" />
        </div>
    </div>
</template>
