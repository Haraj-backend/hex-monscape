package pokestrg

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type extraRole string

const (
	indexExtraRole = "extra_role"

	partnerRole extraRole = "PARTNER"
	enemyRole   extraRole = "ENEMY"
)

type pokemonKey struct {
	ID string `json:"id"`
}

func (k pokemonKey) toDDBKey() map[string]*dynamodb.AttributeValue {
	item, _ := dynamodbattribute.MarshalMap(k)
	return item
}

type pokemonExtraRoleQuery struct {
	ExtraRole extraRole `json:":extra_role"`
}

func (q pokemonExtraRoleQuery) toQueryExpression() *string {
	s := "extra_role = :extra_role"
	return &s
}

func (q pokemonExtraRoleQuery) toQueryExpressionValue() map[string]*dynamodb.AttributeValue {
	attributeValue, _ := dynamodbattribute.MarshalMap(q)
	return attributeValue
}
