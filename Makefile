.PHONY: *

run:
	docker build -t hex-pokebattle -f ./build/package/server/Dockerfile .
	docker run -p 9186:9186 hex-pokebattle

test:
	docker-compose down -v
	env DONT_SEED_POKEMON=true \
		docker-compose up --build --remove-orphans -d

    # waiting for Localstack preparations (DynamoDB tables, etc)
	./wait-localstack.sh -h localhost:4566 -s dynamodb

	env DDB_TABLE_BATTLE_NAME="Battles" \
		DDB_TABLE_GAME_NAME="Games" \
		DDB_TABLE_POKEMON_NAME="Pokemons" \
		LOCALSTACK_ENDPOINT="http://localhost:4566" \
		AWS_REGION=eu-west-1 \
		go test -count=1 ./...
	docker-compose down -v

run-with-ddb:
	docker-compose down -v
	docker-compose up --build --remove-orphans
