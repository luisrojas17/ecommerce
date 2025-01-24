# ecommerce

This repository contains a lambda function.

## Set up project

Follow all of these steps to configure dependencies/libraries:

To create a Go module:
- ** go mod init github.com/luisrojas17/ecommerce

To get driver to connect to database:
- ** go get github.com/go-sql-driver/mysql

To get AWS SDK dependencies:
- **go get github.com/aws/aws-sdk-go-v2/aws
- **go get github.com/aws/aws-sdk-go-v2/config
- **go get github.com/aws/aws-lambda-go/events
- **go get github.com/aws/aws-lambda-go/lambda
- **go get github.com/aws/aws-sdk-go-v2/service/secretsmanager