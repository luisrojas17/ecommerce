package handlers

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
)

// Returns the HTTP status code and message
func Handler(path string, method string, body string, headers map[string]string, request events.APIGatewayV2HTTPRequest) (int, string) {

	fmt.Println("Processing path: " + path + " for method: " + method)

	id := request.PathParameters["id"]

	// Convert the user id and product id. The product id only has to be number.
	idn, _ := strconv.Atoi(id)

	return 400, "Method Invalid"

}
