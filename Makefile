.PHONY: *

run:
	-make down
	docker-compose -f ./deploy/local/deployment/docker-compose.yml up --build

down:
	docker-compose -f ./deploy/local/deployment/docker-compose.yml down --remove-orphans

test:
	-docker-compose -f ./deploy/local/integration-test/docker-compose.yml down --remove-orphans
	docker-compose -f ./deploy/local/integration-test/docker-compose.yml up --build --exit-code-from=integration-test
