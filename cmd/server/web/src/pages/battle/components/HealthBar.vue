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
    <progress id="health" class="w-[400px] h-[32px]" :value="currentHealth" :max="maxHealth"></progress>
</template>