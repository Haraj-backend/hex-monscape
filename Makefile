.PHONY: *

run:
	-make build-frontend
	-make down
	docker-compose -f ./deploy/local/deployment/docker-compose.yml up --build

down:
	docker-compose -f ./deploy/local/deployment/docker-compose.yml down --remove-orphans

test:
	-docker-compose -f ./deploy/local/integration-test/docker-compose.yml down --remove-orphans
	docker-compose -f ./deploy/local/integration-test/docker-compose.yml up --build --exit-code-from=integration-test

build-frontend:
	cd cmd/web; yarn install; yarn build:local
