package shared

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var TestConfig = struct {
	EnvKeyLocalstackEndpoint string
	EnvKeyBattleTableName    string
	EnvKeyGameTableName      string
	EnvKeyPokemonTableName   string
}{
	EnvKeyLocalstackEndpoint: "LOCALSTACK_ENDPOINT",
	EnvKeyBattleTableName:    "DDB_TABLE_BATTLE_NAME",
	EnvKeyGameTableName:      "DDB_TABLE_GAME_NAME",
	EnvKeyPokemonTableName:   "DDB_TABLE_POKEMON_NAME",
}

func NewLocalTestDDBClient() *dynamodb.DynamoDB {
	awsSess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Endpoint: aws.String(os.Getenv(TestConfig.EnvKeyLocalstackEndpoint)),
		},
	}))
	return dynamodb.New(awsSess)
}
