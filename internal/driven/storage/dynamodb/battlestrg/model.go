package battlestrg

import (
	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/Haraj-backend/hex-monscape/internal/driven/storage/dynamodb/shared"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type battleKey struct {
	GameID string `dynamodbav:"game_id"`
}

func (k battleKey) toDDBKey() map[string]*dynamodb.AttributeValue {
	item, _ := dynamodbattribute.MarshalMap(k)
	return item
}

type battleRow struct {
	GameID     string            `dynamodbav:"game_id"`
	State      string            `dynamodbav:"state"`
	Partner    shared.MonsterRow `dynamodbav:"partner"`
	Enemy      shared.MonsterRow `dynamodbav:"enemy"`
	LastDamage lastDamageRow     `dynamodbav:"last_damage"`
}

func toBattleRow(battle entity.Battle) battleRow {
	return battleRow{
		GameID:  battle.GameID,
		State:   string(battle.State),
		Partner: shared.ToMonsterRow(*battle.Partner),
		Enemy:   shared.ToMonsterRow(*battle.Enemy),
		LastDamage: lastDamageRow{
			Partner: battle.LastDamage.Partner,
			Enemy:   battle.LastDamage.Enemy,
		},
	}
}

func (r battleRow) toBattle() *entity.Battle {
	return &entity.Battle{
		GameID:  r.GameID,
		State:   entity.State(r.State),
		Partner: r.Partner.ToMonster(),
		Enemy:   r.Enemy.ToMonster(),
		LastDamage: entity.LastDamage{
			Partner: r.LastDamage.Partner,
			Enemy:   r.LastDamage.Enemy,
		},
	}
}

type lastDamageRow struct {
	Partner int `dynamodbav:"partner"`
	Enemy   int `dynamodbav:"enemy"`
}
