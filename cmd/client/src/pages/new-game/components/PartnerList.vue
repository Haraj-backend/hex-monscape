<script>
import { useStore } from "../../../store"
import PartnerListItem from "./PartnerListItem.vue"

export default {
    props: {
        availablePartners: Array,
    },
    components: {
        PartnerListItem
    },
    setup(props) {
        const store = useStore()

        return {
            availablePartners: props.availablePartners,
            store
        }
    }
}
</script>
<template>
    <div class="partner-list-wrapper flex flex-col gap-y-4 mt-10">
        <partner-list-item
            v-for="partner in availablePartners"
            :key="partner.id"
            :partner="partner"
            :active="store.getActivePartner === partner.id"
        ></partner-list-item>
    </div>
</template>