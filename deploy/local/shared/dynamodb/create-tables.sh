#!/bin/bash

# Create monster table
awslocal dynamodb create-table \
    --table-name monster \
    --attribute-definitions \
        AttributeName=id,AttributeType=S \
        AttributeName=extra_role,AttributeType=S \
    --key-schema \
        AttributeName=id,KeyType=HASH \
    --provisioned-throughput \
        ReadCapacityUnits=1,WriteCapacityUnits=1 \
    --global-secondary-indexes \
        "[
            {
                \"IndexName\": \"extra_role\",
                \"KeySchema\": [{\"AttributeName\":\"extra_role\",\"KeyType\":\"HASH\"}],
                \"Projection\": {
                    \"ProjectionType\": \"ALL\"
                },
                \"ProvisionedThroughput\": {
                    \"ReadCapacityUnits\": 1,
                    \"WriteCapacityUnits\": 1
                }
            }
        ]"

# Create game table
awslocal dynamodb create-table \
    --table-name game \
    --attribute-definitions \
        AttributeName=id,AttributeType=S \
    --key-schema \
        AttributeName=id,KeyType=HASH \
    --provisioned-throughput \
        ReadCapacityUnits=1,WriteCapacityUnits=1

# Create battle table
awslocal dynamodb create-table \
    --table-name battle \
    --attribute-definitions \
        AttributeName=game_id,AttributeType=S \
    --key-schema \
        AttributeName=game_id,KeyType=HASH \
    --provisioned-throughput \
        ReadCapacityUnits=1,WriteCapacityUnits=1
