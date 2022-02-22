<script>
import { computed, onMounted, reactive } from "vue"
import router from "../../routes"
import { useStore } from "../../store"
import PokebattleHTTP from "../../composables/http_client"
import { turnStates } from "../../entity/game"

import HealthBar from './components/HealthBar.vue'
import ControlButton from "./components/ControlButton.vue"

export default {
    components: {
        HealthBar,
        ControlButton
    },
    setup() {
        const store = useStore()
        // show battle base
        const gw = document.getElementById("game-wrapper")
        const currentLocation = window.location.pathname.split("/")
        if (currentLocation.findIndex(loc => loc === 'battle') > -1) {
            gw.style.backgroundImage = `url('https://idev-images-test.s3.eu-west-1.amazonaws.com/hex-pokebattle-backgrounds/battle_base.jpg')`
        }

        // get battle state
        const battleState = computed(() => store.getBattleState)

        const client = new PokebattleHTTP()
        const currentGameID = store.getBattleState.game_id
        const controlState = reactive({
            turn: turnStates.DECIDE_TURN,
            previousTurn: '',
            enemyAttack: false,
            partnerTurn: false,
            message: 'Decide turn...',
        })

        const updateControlState = (battleData) => {
            // update the store
            store.setTheBattle(battleData)
            controlState.previousTurn = controlState.turn
            controlState.turn = battleData.state
            // show enemy attack for random time between 2 to 4 seconds
            if (battleData.state === turnStates.DECIDE_TURN) {
                controlState.enemyAttack = true
                controlState.partnerTurn = false

                controlState.message = controlState.previousTurn === turnStates.DECIDE_TURN
                    ? `Enemy attack!<br/><span class='text-2xl'>You got ${battleData.last_damage.partner} damage</span>`
                    : `Attack the enemy!<br/><span class='text-2xl'>You inflicted ${battleData.last_damage.enemy} damage</span>`
                setTimeout(() => {
                    controlState.enemyAttack = false
                    controlState.message = 'Decide turn...'
                }, 2 * 1000)
            }
            if (battleData.state === turnStates.PARTNER_TURN) {
                controlState.enemyAttack = false
                controlState.partnerTurn = true
                controlState.message = `Your turn`
            }
        }
        // decide turn
        // when the state `DECIDE_TURN` it means we're attacked by the enemy
        // when the state `PARTNER_TURN` then, we should show attack and
        // surrender buttons respectively
        const decideTurn = async () => {
            const resp = await client.decideTurn(currentGameID)
            if (resp.ok) {
                updateControlState(resp.data)
            }
        }

        // attack!
        const attack = async () => {
            const resp = await client.attack(currentGameID)
            if (resp.ok) {
                updateControlState(resp.data)
            }
        }

        // surrender
        const surrender = async () => {
            const resp = await client.surrender(currentGameID)
            if (resp.ok) {
                updateControlState(resp.data)
            }
        }

        // go back
        const exitBattle = () => {
            router.push({ name: 'lounge-screen', params: { state: 'ongoing' } })
        }

        onMounted(() => {
            controlState.turn = battleState.value.state
            if (battleState.value.state === turnStates.PARTNER_TURN) {
                controlState.previousTurn = turnStates.DECIDE_TURN
                controlState.enemyAttack = false
                controlState.partnerTurn = true
                controlState.message = `Your turn`
            }
        })

        return {
            turnStates,
            battleState,
            controlState,
            decideTurn,
            attack,
            surrender,
            exitBattle
        }
    },
    beforeRouteLeave(to, from) {
        if (this.battleState.state !== turnStates.WIN && this.battleState.state !== turnStates.LOSE) {
            const answer = window.confirm('Do you really want to leave? You will lose this battle!')
            if (!answer) return false
        }

        // return back the background if exit battle scene
        if (to.name !== 'battle') {
            const store = useStore()
            const gw = document.getElementById("game-wrapper")
            gw.style.backgroundImage = `url(${store.gameBackground})`
        }
    }
}
</script>
<template>
    <div class="flex flex-col px-16 py-12">
        <!-- Battle scene -->
        <div class="battle-scene-wrapper">
            <!-- Partner -->
            <div class="partner-scene">
                <div class="top-part">
                    <HealthBar
                        :maxHealth="battleState.partner.battle_stats.max_health"
                        :currentHealth="battleState.partner.battle_stats.health"
                    />
                    <p class="partner-name">{{ battleState.partner.name }}</p>
                </div>
                <div class="pokemon-avatar">
                    <img
                        width="256"
                        height="256"
                        :src="battleState.partner.avatar_url"
                        alt="partner_avatar"
                    />
                </div>
            </div>

            <div class="middle-part">
                <div class="battle-description font-bold text-2xl">1st battle</div>
            </div>

            <div class="enemy-scene">
                <div class="top-part">
                    <HealthBar
                        :maxHealth="battleState.enemy.battle_stats.max_health"
                        :currentHealth="battleState.enemy.battle_stats.health"
                    />
                    <p class="enemy-name">{{ battleState.enemy.name }}</p>
                </div>
                <div class="pokemon-avatar">
                    <img
                        width="256"
                        height="256"
                        :src="battleState.enemy.avatar_url"
                        alt="partner_avatar"
                    />
                </div>
            </div>
        </div>

        <!-- Battle control -->
        <div class="battle-control-wrapper">
            <!-- Battle ends, Win or lose status -->
            <template
                v-if="battleState.state === turnStates.WIN || battleState.state === turnStates.LOSE"
            >
                <div class="control-description">
                    <p
                        class="font-bold"
                    >{{ battleState.state === turnStates.WIN ? 'YOU WIN!' : 'YOU LOSE...' }}</p>
                </div>
                <ControlButton @click="exitBattle" type="general" label="Exit" />
            </template>

            <!-- Battle still ongoing -->
            <template v-else>
                <div class="control-description">
                    <p v-html="controlState.message"></p>
                </div>
                <ControlButton
                    v-if="controlState.turn === turnStates.DECIDE_TURN && !controlState.enemyAttack"
                    @click="decideTurn"
                    type="general"
                    label="Continue"
                />
                <div
                    v-else-if="controlState.turn === turnStates.PARTNER_TURN && controlState.partnerTurn"
                    class="flex gap-x-4"
                >
                    <ControlButton @click="attack" type="attack" label="Attack" />
                    <ControlButton @click="surrender" type="surrender" label="Surrender" />
                </div>
            </template>
        </div>
    </div>
</template>

<style scoped>
.battle-scene-wrapper {
    @apply flex justify-between;
}
.battle-control-wrapper {
    @apply flex flex-col w-full items-center justify-center gap-y-4 mt-20;
}
.control-description {
    @apply w-[800px] bg-white py-12;
    @apply shadow-[0_6px_0_rgba(0,0,0,.15)] rounded-lg;
    @apply text-center text-5xl italic;
}
.partner-name,
.enemy-name {
    @apply mt-2 text-xl;
}

.enemy-name {
    @apply text-right;
}

.pokemon-avatar {
    @apply flex justify-center mt-44;
}
</style>