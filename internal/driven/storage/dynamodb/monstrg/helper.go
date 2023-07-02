package monstrg

import (
	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type PokemonSeeder struct {
	partners []entity.Monster
	enemies  []entity.Monster
}

func (s PokemonSeeder) toBatchWriteInput(tableName string) *dynamodb.BatchWriteItemInput {
	writeRequests := make([]*dynamodb.WriteRequest, 0)
	writeRequests = append(writeRequests, pokemonListToWriteRequests(s.partners, partnerRole)...)
	writeRequests = append(writeRequests, pokemonListToWriteRequests(s.enemies, enemyRole)...)

	return &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			tableName: writeRequests,
		},
	}
}

func (s PokemonSeeder) isEmpty() bool {
	return (len(s.enemies) + len(s.partners)) == 0
}

func pokemonListToWriteRequests(pokemons []entity.Monster, role extraRole) []*dynamodb.WriteRequest {
	writeRequests := make([]*dynamodb.WriteRequest, 0)

	for _, pokemon := range pokemons {
		item, _ := dynamodbattribute.MarshalMap(pokemon)
		item["extra_role"], _ = dynamodbattribute.Marshal(string(role))

		putRequest := dynamodb.PutRequest{
			Item: item,
		}

		writeRequest := dynamodb.WriteRequest{
			PutRequest: &putRequest,
		}

		writeRequests = append(writeRequests, &writeRequest)
	}

	return writeRequests
}
