package entity

import (
	"github.com/google/uuid"
	"gopkg.in/validator.v2"
)

type Scenario string

const (
	ScenarioBattle1 Scenario = "BATTLE_1"
	ScenarioBattle2 Scenario = "BATTLE_2"
	ScenarioBattle3 Scenario = "BATTLE_3"
	ScenarioEndGame Scenario = "END_GAME"
)

type Game struct {
	ID         string
	PlayerName string
	Partner    *Monster
	CreatedAt  int64
	BattleWon  int
	Scenario   Scenario
}

type GameConfig struct {
	PlayerName string   `validate:"nonzero"`
	Partner    *Monster `validate:"nonnil"`
	CreatedAt  int64    `validate:"nonzero"`
}

func (c GameConfig) Validate() error {
	return validator.Validate(c)
}

// AdvanceScenario is used for advancing current game scenario. It will
// calculate the next scenario based on game current condition. Beside
// updating game internal scenario into the new one, it also returns the
// new scenario value.
func (g *Game) AdvanceScenario() Scenario {
	// determine next scenario
	nextScenario := g.Scenario
	if g.BattleWon >= 3 {
		nextScenario = ScenarioEndGame
	} else if g.BattleWon == 2 {
		nextScenario = ScenarioBattle3
	} else if g.BattleWon == 1 {
		nextScenario = ScenarioBattle2
	}
	// update internal scenario & return its value
	g.Scenario = nextScenario
	return g.Scenario
}

// IncBattleWon is used for incrementing number of battle won then advancing
// the scenario
func (g *Game) IncBattleWon() {
	g.BattleWon++
	g.AdvanceScenario()
}

func NewGame(cfg GameConfig) (*Game, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}
	g := &Game{
		ID:         uuid.NewString(),
		PlayerName: cfg.PlayerName,
		Partner:    cfg.Partner,
		CreatedAt:  cfg.CreatedAt,
		BattleWon:  0,
		Scenario:   ScenarioBattle1,
	}
	return g, nil
}
