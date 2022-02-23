import { localStorageKeys } from "../composables/constants";
import { getGameScenario, turnStates } from "../entity/game";

// getAppData used to get game and battle state data
const getGameConfig = () => {
  const config = localStorage.getItem(localStorageKeys.GAME_CONFIG);
  if (config) {
    return JSON.parse(config);
  }

  return null;
};

const gameStateCheck = (config) => {
  return config && config.gameData !== null;
};

const battleCheck = (config) => {
  return (
    config &&
    config.battleState !== null &&
    (config.battleState.state === turnStates.DECIDE_TURN ||
      config.battleState.state === turnStates.PARTNER_TURN)
  );
};

export const gameStateMiddleware = (to, from, next) => {
  const gameConfig = getGameConfig();
  if (
    gameStateCheck(gameConfig) &&
    !battleCheck(gameConfig) &&
    to.name !== "lounge-screen"
  ) {
    next({ name: "lounge-screen", params: { state: "ongoing" } });
  } else if (
    gameStateCheck(gameConfig) &&
    battleCheck(gameConfig) &&
    to.name !== "battle"
  ) {
    const scenarioNum = getGameScenario(gameConfig.gameData);
    next({ name: "battle", params: { scenario: `scenario-${scenarioNum}` } });
  } else {
    next();
  }
};
