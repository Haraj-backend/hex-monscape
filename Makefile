.PHONY: *

SEED_POKEMONS=true

run:
	docker build -t hex-pokebattle -f ./build/package/server/Dockerfile .
	docker run -p 9186:9186 hex-pokebattle

test:
	docker-compose down -v
	docker-compose up --build --remove-orphans -d

    # waiting for Localstack preparations (DynamoDB tables, etc)
	sh -c 'sleep 60'

	env DDB_TABLE_BATTLE_NAME="Battles" \
		DDB_TABLE_GAME_NAME="Games" \
		DDB_TABLE_POKEMON_NAME="Pokemons" \
		LOCALSTACK_ENDPOINT="http://localhost:4566" \
		AWS_REGION=eu-west-1 \
		go test -count=1 ./...
	docker-compose down -v

run-ddb:
	docker-compose down -v

	env SEED_POKEMON=true \
		IS_SERVER_MODE=true \
		AWS_ACCESS_KEY_ID=${DEV_AWS_ACCESS_KEY_ID} \
		AWS_SECRET_ACCESS_KEY=${DEV_AWS_SECRET_ACCESS_KEY} \
		AWS_REGION=${REGION_DEV} \
		docker-compose up --build --remove-orphans
	docker-compose down -v
