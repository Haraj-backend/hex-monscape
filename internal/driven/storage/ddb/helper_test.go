package ddb

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	envKeyLocalstackEndpoint = "LOCALSTACK_ENDPOINT"
	envKeyBattleTableName    = "DDB_TABLE_BATTLE_NAME"
	envKeyGameTableName      = "DDB_TABLE_GAME_NAME"
	envKeyMonsterTableName   = "DDB_TABLE_MONSTER_NAME"
)

func newTestDDBClient() *dynamodb.DynamoDB {
	awsSess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Endpoint: aws.String(os.Getenv(envKeyLocalstackEndpoint)),
		},
	}))

	return dynamodb.New(awsSess)
}
