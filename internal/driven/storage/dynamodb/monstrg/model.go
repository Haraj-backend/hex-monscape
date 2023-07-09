package monstrg

import (
	"fmt"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/Haraj-backend/hex-monscape/internal/driven/storage/dynamodb/shared"
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

func toMonsters(items []map[string]*dynamodb.AttributeValue) ([]entity.Monster, error) {
	var monsterRows []shared.MonsterRow
	err := dynamodbattribute.UnmarshalListOfMaps(items, &monsterRows)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal items due to: %w", err)
	}
	monsters := make([]entity.Monster, len(monsterRows))
	for i := 0; i < len(monsterRows); i++ {
		monsters[i] = *monsterRows[i].ToMonster()
	}
	return monsters, nil
}
