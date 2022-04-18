.PHONY: *

TIMESTAMP:=$(shell /bin/date "+%s")
INFRA_STACK_NAME_DEV:=hex-pokebattle-infras
AWS_ACCOUNT_ID:=$(shell aws sts get-caller-identity --query Account --output text)
ECR_REPO_NAME_DEV:=$(shell aws cloudformation describe-stack-resource \
	--stack-name ${INFRA_STACK_NAME_DEV} \
	--logical-resource-id ECRRepoHexPokebattle \
	--query "StackResourceDetail.PhysicalResourceId" --output text)
REMOTE_REPO_DEV:=${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${ECR_REPO_NAME_DEV}

INFRA_STACK_NAME_MYSQL_DEV:=hex-pokebattle-mysql-infras
ECR_MYSQL_REPO_NAME_DEV:=$(shell aws cloudformation describe-stack-resource \
	--stack-name ${INFRA_STACK_NAME_MYSQL_DEV} \
	--logical-resource-id ECRRepoHexPokebattleMysql \
	--query "StackResourceDetail.PhysicalResourceId" --output text)
REMOTE_MYSQL_REPO_DEV:=${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${ECR_MYSQL_REPO_NAME_DEV}

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
		AWS_ACCESS_KEY_ID=awslocal \
		AWS_SECRET_ACCESS_KEY=awslocal \
		TEST_SQL_DSN="root:test1234@tcp(127.0.0.1:3307)/db_pokebattle?timeout=5s" \
		go test -count=1 ./...
	docker-compose down -v

run-with-ddb:
	docker-compose down -v
	docker-compose up --build --remove-orphans

deploy-infras-dev:
	aws cloudformation deploy \
		--region ${AWS_REGION} \
		--template-file ./deploy/aws/infras.yml \
		--stack-name ${INFRA_STACK_NAME_DEV} \
		--capabilities CAPABILITY_NAMED_IAM

build-push-image-dev:
	docker build \
		--build-arg VITE_API_STAGE_PATH=/Dev \
		--build-arg FRONTEND_MODE=lambda \
		-t hex-pokebattle-lambda:latest -f ./build/package/lambda/Dockerfile .
	docker tag hex-pokebattle-lambda:latest ${REMOTE_REPO_DEV}:${TIMESTAMP}

	aws ecr get-login-password | docker login --username AWS --password-stdin ${REMOTE_REPO_DEV}
	docker push ${REMOTE_REPO_DEV}:${TIMESTAMP}

deploy-dev: build-push-image-dev
	sam deploy \
		--region ${AWS_REGION} \
		--stack-name hex-pokebattle \
		--image-repository ${REMOTE_REPO_DEV} \
		--template-file ./deploy/aws/services.yml \
		--capabilities CAPABILITY_NAMED_IAM \
		--parameter-overrides \
			InfraStackName=${INFRA_STACK_NAME_DEV} \
			ImageUri=${REMOTE_REPO_DEV}:${TIMESTAMP}

deploy-infras-dev-mysql:
	aws cloudformation deploy \
		--region ${AWS_REGION} \
		--template-file ./deploy/aws/mysql/infras.yml \
		--stack-name ${INFRA_STACK_NAME_MYSQL_DEV} \
		--capabilities CAPABILITY_NAMED_IAM \
		--parameter-overrides \
			MasterUserPassword=pokebattle1234

build-push-image-dev-mysql:
	docker build \
		--build-arg VITE_API_STAGE_PATH=/Dev \
		--build-arg FRONTEND_MODE=lambda \
		-t hex-pokebattle-lambda-mysql:latest -f ./build/package/lambda-mysql/Dockerfile .
	docker tag hex-pokebattle-lambda-mysql:latest ${REMOTE_MYSQL_REPO_DEV}:${TIMESTAMP}

	aws ecr get-login-password | docker login --username AWS --password-stdin ${REMOTE_MYSQL_REPO_DEV}
	docker push ${REMOTE_MYSQL_REPO_DEV}:${TIMESTAMP}

deploy-dev-mysql: build-push-image-dev-mysql
	sam deploy \
		--region ${AWS_REGION} \
		--stack-name hex-pokebattle-mysql \
		--image-repository ${REMOTE_MYSQL_REPO_DEV} \
		--template-file ./deploy/aws/mysql/services.yml \
		--capabilities CAPABILITY_NAMED_IAM \
		--parameter-overrides \
			InfraStackName=${INFRA_STACK_NAME_MYSQL_DEV} \
			ImageUri=${REMOTE_MYSQL_REPO_DEV}:${TIMESTAMP} \
			MasterUserPassword=pokebattle1234 \
			DatabaseName=db_pokebattle
