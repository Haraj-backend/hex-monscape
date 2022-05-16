<script>
import { computed } from "vue"
import { Icon } from '@iconify/vue';

export default {
    props: {
        partner: Object
    },
    components: {
        Icon
    },
    setup(props) {
        const partner = computed(() => props.partner)
        return {
            partner
        }
    }
}
</script>
<template>
    <div
        v-if="partner"
        class="flex flex-col items-center bg-white w-[400px] h-[490px] py-6 border border-2 border-black rounded-lg"
    >
        <img :src="partner.avatar_url" alt="avatar_url" width="192" height="192" />
        <p class="partner-name text-2xl">{{ partner.name }}</p>
        <p class="text-gray-400">Battle stats</p>
        <div class="flex flex-col w-3/4 mt-6 partner-stats">
            <!-- Health -->
            <div class="flex gap-x-2 items-center w-full px-4 py-2 border-b border-gray-300">
                <div class="stats-property flex items-center w-1/3 gap-x-2">
                    <Icon icon="akar-icons:health" />
                    <span class="property-text">Health</span>
                </div>
                <div class="stats-value font-bold text-green-600">{{ partner.battle_stats.health }}</div>
            </div>

            <!-- Attack -->
            <div class="flex gap-x-2 w-full px-4 py-2 border-b border-gray-300">
                <div class="stats-property flex items-center w-1/3 gap-x-2">
                    <Icon icon="akar-icons:sword" />
                    <span class="property-text">Attack</span>
                </div>
                <div class="stats-value font-bold text-red-600">{{ partner.battle_stats.attack }}</div>
            </div>

            <!-- Defense -->
            <div class="flex gap-x-2 w-full px-4 py-2 border-b border-gray-300">
                <div class="stats-property flex items-center w-1/3 gap-x-2">
                    <Icon icon="akar-icons:shield" />
                    <span class="property-text">Defense</span>
                </div>
                <div class="stats-value font-bold text-blue-600">{{ partner.battle_stats.defense }}</div>
            </div>

            <!-- Speed -->
            <div class="flex gap-x-2 w-full px-4 py-2 border-gray-300">
                <div class="stats-property flex items-center w-1/3 gap-x-2">
                    <Icon icon="wi:strong-wind" />
                    <span class="property-text">Speed</span>
                </div>
                <div class="stats-value font-bold text-orange-600">{{ partner.battle_stats.speed }}</div>
            </div>
        </div>
    </div>
</template>