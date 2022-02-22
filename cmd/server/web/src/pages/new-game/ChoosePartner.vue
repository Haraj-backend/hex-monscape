<script>
import { computed, ref } from "vue"
import { useRouter } from "vue-router"

import { useStore } from "../../store"
import { getAvailablePartners } from "../../entity/partner"
import PartnerList from "./components/PartnerList.vue"
import PartnerCard from "../../components/PartnerCard.vue"
import PokebattleHTTP from "../../composables/http_client"

export default {
    components: {
        PartnerList,
        PartnerCard
    },
    setup() {
        const router = useRouter()
        const store = useStore()
        const client = new PokebattleHTTP()

        let availablePartners = ref([])
        getAvailablePartners().then(p => availablePartners.value = p)

        const chosenPartner = computed(() => availablePartners.value.find(p => p.id === store.getActivePartner))

        const newGame = async () => {
            const r = await client.newGame({
                player_name: store.playerName,
                partner_id: store.partnerData.id
            })
            if (r.ok) {
                store.setTheGame(r.data)
                router.push({ name: 'lounge-screen', params: { state: 'ongoing' } })
            }
        }

        return {
            playerName: store.getPlayerName,
            availablePartners,
            chosenPartner,
            newGame
        }
    }
}
</script>

<template>
    <div v-if="availablePartners.length > 0" class="flex justify-between w-full h-app p-[160px]">
        <div class="left-side">
            <!-- Game title -->
            <h1 class="game-title">NEW GAME</h1>
            <p class="game-description mt-3 text-2xl">
                Choose your partner,
                <span class="player-name">{{ playerName }}</span>
            </p>

            <PartnerList :availablePartners="availablePartners" />

            <!-- Initial actions -->
            <div class="game-initial-actions">
                <button
                    @click="newGame"
                    class="bg-red-600 text-white rounded-lg w-[300px] text-2xl py-2 px-3"
                >Proceed</button>
            </div>
        </div>
        <div class="right-side">
            <PartnerCard :partner="chosenPartner" />
        </div>
    </div>
</template>