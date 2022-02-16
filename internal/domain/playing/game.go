package playing

import "github.com/riandyrn/pokebattle/internal/domain/entity"

type Scenario string

const (
	BATTLE_1 Scenario = "BATTLE_1"
	BATTLE_2 Scenario = "BATTLE_2"
	BATTLE_3 Scenario = "BATTLE_3"
	END_GAME Scenario = "END_GAME"
)

type Game struct {
	ID         string
	PlayerName string
	Partner    entity.Pokemon
	CreatedAt  int64
	BattleWon  int
	Scenario   Scenario
}

func (g *Game) GetNextScenario() Scenario {
	// TODO
	return ""
}

func NewGame(playerName string, partner entity.Pokemon) (*Game, error) {
	// TODO
	return nil, nil
}
