package playing

import (
	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
	"github.com/google/uuid"
	"gopkg.in/validator.v2"
)

type Scenario string

const (
	BATTLE_1 Scenario = "BATTLE_1"
	BATTLE_2 Scenario = "BATTLE_2"
	BATTLE_3 Scenario = "BATTLE_3"
	END_GAME Scenario = "END_GAME"
)

type Game struct {
	ID         string          `json:"id"`
	PlayerName string          `json:"player_name"`
	Partner    *entity.Pokemon `json:"partner"`
	CreatedAt  int64           `json:"created_at"`
	BattleWon  int             `json:"battle_won"`
	Scenario   Scenario        `json:"scenario"`
}

type GameConfig struct {
	PlayerName string          `validate:"nonzero"`
	Partner    *entity.Pokemon `validate:"nonnil"`
	CreatedAt  int64           `validate:"nonzero"`
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
		nextScenario = END_GAME
	} else if g.BattleWon == 2 {
		nextScenario = BATTLE_3
	} else if g.BattleWon == 1 {
		nextScenario = BATTLE_2
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
		Scenario:   BATTLE_1,
	}
	return g, nil
}
