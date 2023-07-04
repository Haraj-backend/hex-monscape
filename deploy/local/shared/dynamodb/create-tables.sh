#!/bin/bash

# Create monster table
awslocal dynamodb create-table \
    --table-name monster \
    --attribute-definitions \
        AttributeName=id,AttributeType=S \
        AttributeName=is_partnerable,AttributeType=N \
    --key-schema \
        AttributeName=id,KeyType=HASH \
    --provisioned-throughput \
        ReadCapacityUnits=1,WriteCapacityUnits=1 \
    --global-secondary-indexes \
        '
            [
                {
                    "IndexName": "is_partnerable",
                    "KeySchema": [{"AttributeName":"is_partnerable","KeyType":"HASH"}],
                    "Projection": {
                        "ProjectionType": "ALL"
                    },
                    "ProvisionedThroughput": {
                        "ReadCapacityUnits": 1,
                        "WriteCapacityUnits": 1
                    }
                }
            ]        
        '

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
