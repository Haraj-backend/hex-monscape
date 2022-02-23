export const turnStates = {
  DECIDE_TURN: "DECIDE_TURN",
  PARTNER_TURN: "PARTNER_TURN",
  WIN: "WIN",
  LOSE: "LOSE",
  END_GAME: "END_GAME",
};

export const getGameScenario = (gameDetails) => {
  return gameDetails.scenario.split("_")[1];
};
