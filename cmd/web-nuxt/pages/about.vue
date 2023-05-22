<template>
    <div id="welcome-screen-wrapper" class="flex relative w-app h-app bg-cover">
        <div class="p-[160px]">
            <!-- Game title -->
            <h1 class="game-title">
                HEX
                <br />POKEBATTLE
            </h1>
            <div class="flex flex-col gap-y-2 mt-4 pr-96 text-lg">
                <p>
                    This is a simple web game implemented using Hexagonal Architecture. It is intended to become Solutions
                    Team reference when they want to implement web app using Hexagonal Architecture. You can see more
                    details
                    <a class="underline" href="https://github.com/Haraj-backend/hex-pokebattle" target="_blank"
                        rel="noopener noreferrer">here.</a>
                </p>
                <p>This web based UI is used to demonstrate how the Hex Pokebattle game would work.</p>
            </div>

            <!-- Initial actions -->
            <div class="game-initial-actions">
                <NuxtLink to="/" class="button w-[264px]">
                    Back</NuxtLink>
            </div>
        </div>
    </div>
</template>