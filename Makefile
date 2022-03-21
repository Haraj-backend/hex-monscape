.PHONY: *

run:
	docker build -t hex-pokebattle -f ./build/package/server/Dockerfile .
	docker run -p 9186:9186 hex-pokebattle

test:
	docker-compose down -v
	env DONT_SEED_POKEMON=true \
		docker-compose up --build --remove-orphans -d

    # waiting for Localstack preparations (DynamoDB tables, etc)
	./deploy/local/wait-localstack.sh -h localhost:4566 -s dynamodb

	env DDB_TABLE_BATTLE_NAME="Battles" \
		DDB_TABLE_GAME_NAME="Games" \
		DDB_TABLE_POKEMON_NAME="Pokemons" \
		LOCALSTACK_ENDPOINT="http://localhost:4566" \
		AWS_REGION=eu-west-1 \
		go test -count=1 ./...
	docker-compose down -v

build-local:
	sam build \
			--template-file ./build/package/lambda/template.yml \
			--parameter-overrides \
				VitePokebattleUrl=http://localhost:3000 \
				LocalDeploymentEnabled=true \
				LocalDeploymentEndpoint=http://host.docker.internal:4566

deploy-local:
	docker-compose down -v
	docker-compose up --build --remove-orphans -d
	make build-local
	./deploy/local/wait-localstack.sh -h localhost:4566 -s dynamodb

	sam local start-api \
		--warm-containers LAZY \
		--parameter-overrides \
			LocalDeploymentEnabled=true \
			LocalDeploymentEndpoint=http://host.docker.internal:4566
