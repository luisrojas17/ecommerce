package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {

	// Pass the request to handler, in this case Execute method is the init point
	// to start to execute the lambda function.
	lambda.Start(Execute)

}

func Execute(ctx context.Context, request events.APIGatewayV2HTTPRequest) (*events.APIGatewayProxyResponse, error) {

	// Init AWS configuration for lambda
	//aws.Init()

}
