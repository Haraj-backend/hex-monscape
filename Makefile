.PHONY: *

run:
	docker build -t hex-pokebattle -f ./build/package/server/Dockerfile .
	docker run -p 9186:9186 hex-pokebattle

test:
	docker-compose down -v
	docker-compose up --build --remove-orphans -d

    # waiting for Localstack preparations (DynamoDB tables, etc)
	sh -c 'sleep 60'

	env DDB_TABLE_BATTLE_NAME="Battles" \
		LOCALSTACK_ENDPOINT="http://localhost:4566" \
		AWS_REGION=eu-west-1 \
		go test -count=1 ./...
	docker-compose down -v
