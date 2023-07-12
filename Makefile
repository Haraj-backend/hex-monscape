.PHONY: *

# by default, run the rest memory server variant
run:
	make run-rest-memory

# run the rest memory server variant
run-rest-memory:
	-docker compose -f ./deploy/local/run/rest-memory/docker-compose.yml down --remove-orphans
	docker compose -f ./deploy/local/run/rest-memory/docker-compose.yml up --build --attach=server --attach=client

# run the rest dynamodb server variant
run-rest-dynamodb:
	-docker compose -f ./deploy/local/run/rest-dynamodb/docker-compose.yml down --remove-orphans
	docker compose -f ./deploy/local/run/rest-dynamodb/docker-compose.yml up --build --attach=server --attach=client

# run the rest mysql server variant
run-rest-mysql:
	-docker compose -f ./deploy/local/run/rest-mysql/docker-compose.yml down --remove-orphans
	docker compose -f ./deploy/local/run/rest-mysql/docker-compose.yml up --build --attach=server --attach=client

# run the lambda server variant
run-lambda:
	-docker compose -f ./deploy/local/run/lambda/docker-compose.yml down --remove-orphans
	docker compose -f ./deploy/local/run/lambda/docker-compose.yml up --build --attach=lambda

# execute the tests
test:
	-docker compose -f ./deploy/local/tests/docker-compose.yml down --remove-orphans
	docker compose -f ./deploy/local/tests/docker-compose.yml up --build --attach=tests --exit-code-from=tests

# building all services to ensure everything is okay
test-build-all:
	docker build -t hex-monscape-client -f ./build/package/client/Dockerfile .
	docker build -t hex-monscape-server -f ./build/package/server/Dockerfile .
	docker build -t hex-monscape-lambda -f ./build/package/lambda/Dockerfile .
	docker build -t hex-monscape-lambda-local -f ./build/package/lambda/Dockerfile.local .