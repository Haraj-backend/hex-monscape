package monstrg

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

type monsterKey struct {
	ID string `json:"id"`
}

func (k monsterKey) toDDBKey() map[string]*dynamodb.AttributeValue {
	item, _ := dynamodbattribute.MarshalMap(k)
	return item
}

type monsterExtraRoleQuery struct {
	ExtraRole extraRole `json:":extra_role"`
}

func (q monsterExtraRoleQuery) toQueryExpression() *string {
	s := "extra_role = :extra_role"
	return &s
}

func (q monsterExtraRoleQuery) toQueryExpressionValue() map[string]*dynamodb.AttributeValue {
	attributeValue, _ := dynamodbattribute.MarshalMap(q)
	return attributeValue
}
