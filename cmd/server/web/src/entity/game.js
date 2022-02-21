export const turnStates = {
  DECIDE_TURN: "DECIDE_TURN",
  PARTNER_TURN: "PARTNER_TURN",
  WIN: "WIN",
  LOSE: "LOSE",
};

export const getGameScenario = (gameDetails) => {
  return gameDetails.data.scenario.split("_")[1];
};
