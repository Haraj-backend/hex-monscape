package gamestrg

import (
	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/Haraj-backend/hex-monscape/internal/driven/storage/dynamodb/shared"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type gameKey struct {
	ID string `json:"id"`
}

func (k gameKey) toDDBKey() map[string]*dynamodb.AttributeValue {
	item, _ := dynamodbattribute.MarshalMap(k)
	return item
}

type gameRow struct {
	ID         string            `dynamodbav:"id"`
	PlayerName string            `dynamodbav:"player_name"`
	Partner    shared.MonsterRow `dynamodbav:"partner"`
	CreatedAt  int64             `dynamodbav:"created_at"`
	BattleWon  int               `dynamodbav:"battle_won"`
	Scenario   string            `dynamodbav:"scenario"`
}

func toGameRow(game entity.Game) gameRow {
	return gameRow{
		ID:         game.ID,
		PlayerName: game.PlayerName,
		Partner:    shared.ToMonsterRow(*game.Partner),
		CreatedAt:  game.CreatedAt,
		BattleWon:  game.BattleWon,
		Scenario:   string(game.Scenario),
	}
}

func (r gameRow) toGame() *entity.Game {
	return &entity.Game{
		ID:         r.ID,
		PlayerName: r.PlayerName,
		Partner:    r.Partner.ToMonster(),
		CreatedAt:  r.CreatedAt,
		BattleWon:  r.BattleWon,
		Scenario:   entity.Scenario(r.Scenario),
	}
}
