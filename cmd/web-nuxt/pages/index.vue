<template>
    <div id="welcome-screen-wrapper" class="flex relative w-app h-app bg-cover">
        <div class="p-[160px]">
            <!-- Game title -->
            <h1 class="game-title">
                HEX
                <br />POKEBATTLE
            </h1>

            <!-- Initial actions -->
            <div class="game-initial-actions">

                <NuxtLink to="/new-game/player-name" class="button font-bold w-[264px]">
                    New Game
                </NuxtLink>
                <NuxtLink to="/about" class="button w-[264px]">
                    About
                </NuxtLink>
            </div>
        </div>
    </div>
</template>