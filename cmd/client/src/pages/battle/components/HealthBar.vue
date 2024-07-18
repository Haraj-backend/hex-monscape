<script>
import { reactive } from "vue"

export default {
    props: {
        maxHealth: Number,
        currentHealth: Number
    },
}
</script>

<template>
    <progress id="health" class="w-full h-[32px]" :value="currentHealth" :max="maxHealth"></progress>
</template>

<style scoped>
@media (min-width: 1060px) {
    #health {
        @apply w-[400px];
    }
}
</style>