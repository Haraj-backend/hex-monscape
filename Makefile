.PHONY: *

run:
	docker build -t hex-pokebattle -f ./build/package/server/Dockerfile .
	docker run -p 9186:9186 hex-pokebattle

test:
	docker-compose down -v
	docker-compose up --build --remove-orphans -d

    # waiting for Localstack preparations (DynamoDB tables, etc)
	sh -c 'sleep 10'

	env HARAJ_POKEBATTLE_AWS_REGION="ap-southeast-1" \
		HARAJ_POKEBATTLE_AWS_DYNAMODB_URL="http://localhost:4566" \
		go test -count=1 ./...
	docker-compose down -v