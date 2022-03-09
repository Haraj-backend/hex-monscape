package pokestrg

import (
	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type PokemonSeeder struct {
	partners []entity.Pokemon
	enemies  []entity.Pokemon
}

func (s PokemonSeeder) toBatchWriteInput(tableName string) *dynamodb.BatchWriteItemInput {
	writeRequests := make([]*dynamodb.WriteRequest, 0)
	writeRequests = append(writeRequests, pokemonListToWriteRequests(s.partners, PARTNER)...)
	writeRequests = append(writeRequests, pokemonListToWriteRequests(s.enemies, ENEMY)...)

	return &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			tableName: writeRequests,
		},
	}
}

func pokemonListToWriteRequests(pokemons []entity.Pokemon, role ExtraRole) []*dynamodb.WriteRequest {
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
