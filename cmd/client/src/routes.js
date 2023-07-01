import * as VueRouter from "vue-router";
import WelcomeScreen from "./pages/welcome-screen/WelcomeScreen.vue";
import AboutScreen from "./pages/welcome-screen/AboutScreen.vue";
import LoungeScreen from "./pages/lounge-screen/LoungeScreen.vue";
import NGPlayerName from "./pages/new-game/PlayerName.vue";
import NGChoosePartner from "./pages/new-game/ChoosePartner.vue";
import BattleScene from "./pages/battle/BattleScene.vue";
import { gameStateMiddleware } from "./middlewares/game";

const routes = [
  {
    path: "/",
    name: "welcome-screen",
    component: WelcomeScreen,
    beforeEnter: gameStateMiddleware,
  },
  {
    path: "/about",
    name: "about-screen",
    component: AboutScreen,
    beforeEnter: gameStateMiddleware,
  },
  {
    path: "/new-game/player-name",
    name: "player-name",
    component: NGPlayerName,
    beforeEnter: gameStateMiddleware,
  },
  {
    path: "/new-game/choose-partner",
    name: "choose-partner",
    component: NGChoosePartner,
    beforeEnter: gameStateMiddleware,
  },
  {
    path: "/lounge/:state",
    name: "lounge-screen",
    component: LoungeScreen,
    beforeEnter: gameStateMiddleware,
  },
  {
    path: "/battle/:scenario",
    name: "battle",
    component: BattleScene,
    beforeEnter: gameStateMiddleware,
  },
];

const router = VueRouter.createRouter({
  history: VueRouter.createWebHistory(import.meta.env.VITE_API_STAGE_PATH),
  routes,
});

export default router;
