.PHONY: *

TIMESTAMP:=$(shell date +%s)
AWS_ACCOUNT_ID=$(shell aws sts get-caller-identity --query Account --output text)
REMOTE_REPO=${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/hex-pokebattle-dev

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

run-with-ddb:
	docker-compose down -v
	docker-compose up --build --remove-orphans

build-push-image-dev:
	docker build \
		--build-arg VITE_API_STAGE_PATH=/Dev \
		--build-arg FRONTEND_MODE=lambda-dev \
		-t hex-pokebattle-lambda:latest -f ./build/package/lambda/Dockerfile .
	docker tag hex-pokebattle-lambda:latest ${REMOTE_REPO}:${TIMESTAMP}

	aws ecr get-login-password | docker login --username AWS --password-stdin ${REMOTE_REPO}
	docker push ${REMOTE_REPO}:${TIMESTAMP}

deploy-dev:
	make build-push-image-dev
	sam deploy \
		--region ${AWS_REGION} \
		--stack-name hex-pokebattle \
		--image-repository ${REMOTE_REPO} \
		--template-file ./deploy/aws/deploy.yml \
		--capabilities CAPABILITY_NAMED_IAM \
		--parameter-overrides \
			StageName=Dev \
			LocalDeploymentEnabled=false \
			ImageUri=${REMOTE_REPO}:${TIMESTAMP}
