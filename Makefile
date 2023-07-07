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

# execute the tests
test:
	-docker compose -f ./deploy/local/tests/docker-compose.yml down --remove-orphans
	docker compose -f ./deploy/local/tests/docker-compose.yml up --build --attach=tests --exit-code-from=tests
