package ddb

import (
	"fmt"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type monsterKey struct {
	ID string `json:"id"`
}

func (k monsterKey) toDDBKey() map[string]*dynamodb.AttributeValue {
	item, _ := dynamodbattribute.MarshalMap(k)
	return item
}

type monsterRow struct {
	ID            string         `dynamodbav:"id"`
	Name          string         `dynamodbav:"name"`
	BattleStats   battleStatsRow `dynamodbav:"battle_stats"`
	AvatarURL     string         `dynamodbav:"avatar_url"`
	IsPartnerable int            `dynamodbav:"is_partnerable"`
}

type battleStatsRow struct {
	Health    int `dynamodbav:"health"`
	MaxHealth int `dynamodbav:"max_health"`
	Attack    int `dynamodbav:"attack"`
	Defense   int `dynamodbav:"defense"`
	Speed     int `dynamodbav:"speed"`
}

func newMonsterRow(m entity.Monster) monsterRow {
	return monsterRow{
		ID:   m.ID,
		Name: m.Name,
		BattleStats: battleStatsRow{
			Health:    m.BattleStats.Health,
			MaxHealth: m.BattleStats.MaxHealth,
			Attack:    m.BattleStats.Attack,
			Defense:   m.BattleStats.Defense,
			Speed:     m.BattleStats.Speed,
		},
		AvatarURL: m.AvatarURL,
	}
}

func (r monsterRow) toMonster() *entity.Monster {
	return &entity.Monster{
		ID:   r.ID,
		Name: r.Name,
		BattleStats: entity.BattleStats{
			Health:    r.BattleStats.Health,
			MaxHealth: r.BattleStats.MaxHealth,
			Attack:    r.BattleStats.Attack,
			Defense:   r.BattleStats.Defense,
			Speed:     r.BattleStats.Speed,
		},
		AvatarURL: r.AvatarURL,
	}
}

func toMonsters(items []map[string]*dynamodb.AttributeValue) ([]entity.Monster, error) {
	var rows []monsterRow
	err := dynamodbattribute.UnmarshalListOfMaps(items, &rows)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal items due to: %w", err)
	}

	monsters := make([]entity.Monster, len(rows))
	for i := 0; i < len(rows); i++ {
		monsters[i] = *rows[i].toMonster()
	}

	return monsters, nil
}

type gameKey struct {
	ID string `json:"id"`
}

func (k gameKey) toDDBKey() map[string]*dynamodb.AttributeValue {
	item, _ := dynamodbattribute.MarshalMap(k)
	return item
}

type gameRow struct {
	ID         string     `dynamodbav:"id"`
	PlayerName string     `dynamodbav:"player_name"`
	Partner    monsterRow `dynamodbav:"partner"`
	CreatedAt  int64      `dynamodbav:"created_at"`
	BattleWon  int        `dynamodbav:"battle_won"`
	Scenario   string     `dynamodbav:"scenario"`
}

func newGameRow(game entity.Game) gameRow {
	return gameRow{
		ID:         game.ID,
		PlayerName: game.PlayerName,
		Partner:    newMonsterRow(*game.Partner),
		CreatedAt:  game.CreatedAt,
		BattleWon:  game.BattleWon,
		Scenario:   string(game.Scenario),
	}
}

func (r gameRow) toGame() *entity.Game {
	return &entity.Game{
		ID:         r.ID,
		PlayerName: r.PlayerName,
		Partner:    r.Partner.toMonster(),
		CreatedAt:  r.CreatedAt,
		BattleWon:  r.BattleWon,
		Scenario:   entity.Scenario(r.Scenario),
	}
}

type battleKey struct {
	GameID string `dynamodbav:"game_id"`
}

func (k battleKey) toDDBKey() map[string]*dynamodb.AttributeValue {
	item, _ := dynamodbattribute.MarshalMap(k)
	return item
}

type battleRow struct {
	GameID     string        `dynamodbav:"game_id"`
	State      string        `dynamodbav:"state"`
	Partner    monsterRow    `dynamodbav:"partner"`
	Enemy      monsterRow    `dynamodbav:"enemy"`
	LastDamage lastDamageRow `dynamodbav:"last_damage"`
}

type lastDamageRow struct {
	Partner int `dynamodbav:"partner"`
	Enemy   int `dynamodbav:"enemy"`
}

func newBattleRow(battle entity.Battle) battleRow {
	return battleRow{
		GameID:  battle.GameID,
		State:   string(battle.State),
		Partner: newMonsterRow(*battle.Partner),
		Enemy:   newMonsterRow(*battle.Enemy),
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
		Partner: r.Partner.toMonster(),
		Enemy:   r.Enemy.toMonster(),
		LastDamage: entity.LastDamage{
			Partner: r.LastDamage.Partner,
			Enemy:   r.LastDamage.Enemy,
		},
	}
}
