echo "Local Deployment"

# Create dynamodb table for Pokemons
awslocal dynamodb create-table \
    --table-name Pokemons \
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
                \"IndexName\": \"index_extra_role\",
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

# Create dynamodb table for Games
awslocal dynamodb create-table \
    --table-name Games \
    --attribute-definitions \
        AttributeName=id,AttributeType=S \
    --key-schema \
        AttributeName=id,KeyType=HASH \
    --provisioned-throughput \
        ReadCapacityUnits=1,WriteCapacityUnits=1

# Create dynamodb table for Battles
awslocal dynamodb create-table \
    --table-name Battles \
    --attribute-definitions \
        AttributeName=game_id,AttributeType=S \
    --key-schema \
        AttributeName=game_id,KeyType=HASH \
    --provisioned-throughput \
        ReadCapacityUnits=1,WriteCapacityUnits=1
