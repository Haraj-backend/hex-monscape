import * as VueRouter from "vue-router";
import WelcomeScreen from "./pages/welcome-screen/WelcomeScreen.vue";
import LoungeScreen from "./pages/lounge-screen/LoungeScreen.vue";
import NGPlayerName from "./pages/new-game/PlayerName.vue";
import NGChoosePartner from "./pages/new-game/ChoosePartner.vue";
import BattleScene from "./pages/battle/BattleScene.vue";

const routes = [
  {
    path: "/",
    name: "welcome-screen",
    component: WelcomeScreen,
  },
  {
    path: "/new-game/player-name",
    name: "player-name",
    component: NGPlayerName,
  },
  {
    path: "/new-game/choose-partner",
    name: "choose-partner",
    component: NGChoosePartner,
  },
  {
    path: "/lounge/:state",
    name: "lounge-screen",
    component: LoungeScreen,
  },
  {
    path: "/battle/:scenario",
    name: "battle",
    component: BattleScene,
  },
];

const router = VueRouter.createRouter({
  history: VueRouter.createWebHistory(),
  routes,
});

export default router;
