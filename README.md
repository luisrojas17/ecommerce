# ecommerce

This repository contains a lambda function which process and manage all requests from API ecommerce to Database.

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

## Dependencies
See the project ecommerceUser.

Additionally, you must configure next: 
- Lambda function in order to manage all request from API ecommerce to database. 
- Two API gateway.
    - First: API Gateway to manage all operations to communicate with our ecommerce-app-jlrg S3 bucket.
    - Second: API Gateway (API ecommerce) to handled all requests from Front-End to Back-End. This API must have:
        - An authorizer in order to request token (JWT) to authenticate each request. The autorizer is: "ecommerce-authorizer" and it is mandatory.
        - A lambda integration because the lambda will be executing all operations to database. The lambda is: "ecommerce"

## Compilation
This project has to be compiled in linux format since AWS lambda runtime for Go only 
recognize lambda Linux executable. So you have to set next environment variables:

- set GOOS=linux
- set GOARCH=amd64

See the scripts "compile.bat" or "compile.sh"

## AWS resources configuration

### Create Lambda

When you are going to load the code zip file for your lambda function you will have to add next  environment varirables:

- SecretName=ecommerce_secret
- UrlPrefix=/v1/ecommerce/ 

Add permissions to lambda to connect to secretmanager in order to read secretes to connect to database.
    - Go to lambda "ecommerce"
    - Go to Configuration -> Permissions -> Role name: ecommerce-role-q9tdulvb 
    