<script>
export default {
    props: {
        type: String,
        label: String,
    }
}
</script>
<template>
    <button :class="`control-button control-${type}`">{{ label }}</button>
</template>

<style scoped>
.control-button {
    @apply w-[400px] py-4;
    @apply shadow-[0_4px_0_rgba(0,0,0,.2)] rounded-lg;
    @apply text-center font-bold;
}
.control-attack {
    @apply bg-red-500;
    @apply text-white;
}
.control-surrender {
    @apply bg-white;
}
.control-general {
    @apply bg-orange-500;
}

@media (max-width: 1060px) {
    .control-button {
        @apply w-full block mb-4;
    }
}
</style>