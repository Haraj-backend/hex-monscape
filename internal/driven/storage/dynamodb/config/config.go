package config

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"gopkg.in/validator.v2"
)

const (
	envAwsRegion      = "HARAJ_POKEBATTLE_AWS_REGION"
	envAwsDynamoDBUrl = "HARAJ_POKEBATTLE_AWS_DYNAMODB_URL"
)

type DynamoDBStorageConfig struct {
	DB dynamodbiface.DynamoDBAPI `validate:"nonnil"`
}

func (c DynamoDBStorageConfig) Validate() error {
	return validator.Validate(c)
}

func BuildDynamoDBInstance() dynamodbiface.DynamoDBAPI {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String(os.Getenv(envAwsRegion)),
		Endpoint: aws.String(os.Getenv(envAwsDynamoDBUrl)),
	}))

	return dynamodb.New(sess)
}
