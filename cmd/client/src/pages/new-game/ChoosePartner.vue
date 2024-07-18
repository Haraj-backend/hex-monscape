<script>
import { computed, ref } from "vue"
import { useRouter } from "vue-router"

import { useStore } from "../../store"
import { getAvailablePartners } from "../../entity/partner"
import PartnerList from "./components/PartnerList.vue"
import PartnerCard from "../../components/PartnerCard.vue"
import MonscapeHTTP from "../../composables/http_client"

export default {
    components: {
        PartnerList,
        PartnerCard
    },
    setup() {
        const router = useRouter()
        const store = useStore()
        const client = new MonscapeHTTP()

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
    <div v-if="availablePartners.length > 0" class="flex flex-col justify-between p-10 pt-[160px] sm:p-[160px]">
        <div class="">
            <!-- Game title -->
            <h1 class="game-title">NEW GAME</h1>
            <p class="game-description mt-3 text-2xl">
                Choose your partner,
                <span class="player-name">{{ playerName }}</span>
            </p>
        </div>
        <div class="flex flex-col-reverse lg:flex-row justify-between">
            <div class="left-side">
                <PartnerList :availablePartners="availablePartners" />
    
                <!-- Initial actions -->
                <div class="game-initial-actions">
                    <button
                        @click="newGame"
                        class="bg-red-600 text-white rounded-lg w-full lg:w-[300px] text-2xl py-2 px-3"
                    >Proceed</button>
                </div>
            </div>
            <div class="right-side mt-4 lg:mt-0">
                <PartnerCard :partner="chosenPartner" />
            </div>
        </div>
    </div>
</template>