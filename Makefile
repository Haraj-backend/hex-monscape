.PHONY: *

run:
	-make down
	docker-compose -f ./deploy/local/run/docker-compose.yml up --build

down:
	docker-compose -f ./deploy/local/run/docker-compose.yml down --remove-orphans

test:
	-docker-compose -f ./deploy/local/tests/docker-compose.yml down --remove-orphans
	docker-compose -f ./deploy/local/tests/docker-compose.yml up --build --exit-code-from=exec-tests
