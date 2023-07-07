.PHONY: *

run:
	make run-rest-memory

run-rest-memory:
	-docker compose -f ./deploy/local/run/rest-memory/docker-compose.yml down --remove-orphans
	docker compose -f ./deploy/local/run/rest-memory/docker-compose.yml up --build --attach=server --attach=client

run-rest-dynamodb:
	-docker compose -f ./deploy/local/run/rest-dynamodb/docker-compose.yml down --remove-orphans
	docker compose -f ./deploy/local/run/rest-dynamodb/docker-compose.yml up --build --attach=server --attach=client

run-rest-mysql:
	-docker compose -f ./deploy/local/run/rest-mysql/docker-compose.yml down --remove-orphans
	docker compose -f ./deploy/local/run/rest-mysql/docker-compose.yml up --build --attach=server --attach=client

test:
	-docker compose -f ./deploy/local/tests/docker-compose.yml down --remove-orphans
	docker compose -f ./deploy/local/tests/docker-compose.yml up --build --attach=tests --exit-code-from=tests
