package main

import (
	"context"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/luisrojas17/ecommerce/aws"
	"github.com/luisrojas17/ecommerce/db"
	"github.com/luisrojas17/ecommerce/handlers"
)

func main() {

	// Pass the request to handler, in this case Execute method is the init point
	// to start to execute the lambda function.
	lambda.Start(Execute)

}

// This lambda function is invoke by API Gateway.
func Execute(ctx context.Context, request events.APIGatewayV2HTTPRequest) (*events.APIGatewayProxyResponse, error) {

	// Init AWS configuration for lambda
	aws.Init()

	if !checkEnvironmentParameters() {
		panic("Error. SecretName and UrlPrefix were not specified.")
	}

	// This is a variable definition for pointer
	var response *events.APIGatewayProxyResponse
	prefix := os.Getenv("UrlPrefix")

	// RawPath contains all URL called.
	// Replace all strings in the URL according to prefix by empty string. -1 means all strings that matching
	path := strings.Replace(request.RawPath, prefix, "", -1)

	// Get HTTP method
	method := request.RequestContext.HTTP.Method

	// Get hte body to process it
	body := request.Body

	// Ge the headers
	header := request.Headers

	// Connecto to database
	db.Connect()

	//
	status, message := handlers.Handler(path, method, body, header, request)

	// We specify headers map for response
	headersToResponse := map[string]string{
		"Content-Type": "application/json",
	}

	// This is the use of a pointer
	response = &events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       string(message),
		Headers:    headersToResponse,
	}

	return response, nil
}

/*
This function check if the secret name was prvided into environment.
*/
func checkEnvironmentParameters() bool {
	var exist bool

	_, exist = os.LookupEnv("SecretName")
	if !exist {
		return exist
	}

	_, exist = os.LookupEnv("UrlPrefix")
	if !exist {
		return exist
	}

	return exist
}
