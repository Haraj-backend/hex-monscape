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
	-docker-compose -f ./deploy/integration_test/docker-compose.yml down --remove-orphans
	docker-compose -f ./deploy/integration_test/docker-compose.yml up --build --exit-code-from=integration_test

run-with-ddb:
	docker-compose down -v
	docker-compose up --build --remove-orphans

deploy-infras-dev-ddb:
	aws cloudformation deploy \
		--region ${AWS_REGION} \
		--template-file ./deploy/aws/infras.yml \
		--stack-name ${INFRA_STACK_NAME_DEV} \
		--capabilities CAPABILITY_NAMED_IAM

build-push-image-dev-ddb:
	docker build \
		--build-arg VITE_API_STAGE_PATH=/Dev \
		--build-arg FRONTEND_MODE=lambda \
		-t hex-pokebattle-lambda:latest -f ./build/package/lambda/Dockerfile .
	docker tag hex-pokebattle-lambda:latest ${REMOTE_REPO_DEV}:${TIMESTAMP}

	aws ecr get-login-password | docker login --username AWS --password-stdin ${REMOTE_REPO_DEV}
	docker push ${REMOTE_REPO_DEV}:${TIMESTAMP}

deploy-dev-ddb: build-push-image-dev-ddb
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
