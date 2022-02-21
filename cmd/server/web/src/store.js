import { defineStore } from "pinia";
import { localStorageKeys } from "./composables/constants";

export const useStore = defineStore("main", {
  state: () => ({
    bg: {
      url: "",
      d: "",
    },
    playerName: "",
    chosenPartner: "pikachu",
    partnerData: null,
    gameData: null,
    battleState: null,
  }),
  getters: {
    gameBackground(state) {
      return state.bg.url;
    },

    getPlayerName(state) {
      return state.playerName;
    },

    getActivePartner(state) {
      return state.chosenPartner;
    },

    getPartnerData(state) {
      return state.partnerData;
    },

    getGameData(state) {
      return state.gameData;
    },

    getBattleState(state) {
      return state.battleState;
    },
  },

  actions: {
    updateLS(gameConfig) {
      // sync. the state with config from storage
      if (gameConfig) {
        this.$state = Object.assign(this.$state, JSON.parse(gameConfig));
      }
      localStorage.setItem(
        localStorageKeys.GAME_CONFIG,
        JSON.stringify(this.$state)
      );
    },

    setGameBackground(newURL, todayDate) {
      const gameCfg = localStorage.getItem(localStorageKeys.GAME_CONFIG);
      // if config not exist on storage,
      // just set by the parameters and return
      if (!gameCfg) {
        this.bg.url = newURL;
        this.bg.d = todayDate;
        return;
      }

      const { bg } = JSON.parse(gameCfg);
      // set the state
      this.bg.url = bg.url;
      this.bg.d = bg.d;

      if (todayDate !== bg.d) {
        this.bg.url = newURL;
        this.bg.d = todayDate;
      }

      // as this method called by root template,
      // all configs on storage will automatically
      // set to store's state
      this.updateLS(gameCfg);
    },

    setTheGame(gameData) {
      this.gameData = Object.assign(
        {},
        {
          id: gameData.id,
          createdAt: gameData.created_at,
          battleWon: gameData.battle_won,
          scenario: gameData.scenario,
        }
      );
      this.updateLS();
    },

    setTheBattle(battleData) {
      this.battleState = Object.assign({}, battleData);
      this.updateLS();
    },

    setPlayerName(name) {
      this.playerName = name;
      this.updateLS();
    },

    choosePartner(partner) {
      this.chosenPartner = partner.id;
      this.partnerData = Object.assign({}, partner);
      this.updateLS();
    },
  },
});
