package battlestrg

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type battleKey struct {
	GameID string `json:"game_id"`
}

func (k battleKey) toDDBKey() map[string]*dynamodb.AttributeValue {
	item, _ := dynamodbattribute.MarshalMap(k)
	return item
}
