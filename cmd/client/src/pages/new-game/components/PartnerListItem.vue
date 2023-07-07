<script>
import { toRef } from "vue"
import { useStore } from "../../../store"
import { Icon } from '@iconify/vue';
import { storeToRefs } from "pinia";

export default {
    props: {
        active: Boolean,
        partner: Object
    },
    components: {
        Icon
    },
    setup(props) {
        const store = useStore()
        const { id, name, icon, color } = props.partner
        const choosePartner = () => store.choosePartner(props.partner)

        return {
            id,
            name,
            icon,
            color,
            choosePartner
        }
    }
}
</script>

<template>
    <div
        class="flex items-center gap-x-2 w-[300px] py-2 pl-3 bg-white rounded-lg border border-2"
        :class="active ? 'border-black' : 'border-[rgba(0,0,0,.2)]'"
        @click="choosePartner"
    >
        <Icon width="40" height="40" :icon="icon" :color="active ? color : '#000'" />
        <p class="partner-name">{{ name }}</p>
    </div>
</template>