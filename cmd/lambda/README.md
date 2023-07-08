# Online Game Demo

This is the application we use to deploy the online game demo. Essentially it runs on top of [AWS Lambda](https://aws.amazon.com/lambda/) & [AWS API Gateway](https://aws.amazon.com/api-gateway/). As for the storage it uses [AWS DynamoDB](https://aws.amazon.com/dynamodb/).

You can see how we deploy this application using [this workflow](../.github/workflows/deploy.yml#L51-L113). All of changes created by this workflow will be reflected in https://hex-monscape.haraj.app.