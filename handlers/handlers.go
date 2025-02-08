package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"

	"github.com/luisrojas17/ecommerce/auth"
	"github.com/luisrojas17/ecommerce/routers"
)

const METHOD_INVALID = "Method Invalid."

// Returns the HTTP status code and message to describe the error.
func Handler(path string, method string, body string, headers map[string]string, request events.APIGatewayV2HTTPRequest) (int, string) {

	fmt.Println("Processing request from path: " + path + " for method: " + method)

	id := request.PathParameters["id"]

	// Convert the user id and product id. The product id only has to be number.
	idn, _ := strconv.Atoi(id)

	// Check if token is valid and it had not expired.
	// Make sure that the user is using the access_token in order to the validation can get the Username attribute.
	ok, statusCode, user := authorize(path, method, headers)

	if !ok {
		return statusCode, user
	}

	// To Do: process path in order to split from token "/" instead of use UrlPrefix variable.
	operation := strings.Split(path, "/")

	// We always get the first item from the array because that item contains the path
	// to identify the operation to do
	switch operation[0] {

	case "user":
		return Users(body, path, method, user, id, request)

	case "product":
		return Products(body, path, method, user, idn, request)

	case "stock":
		return Stock(body, path, method, user, idn, request)

	case "address":
		return Address(body, path, method, user, idn, request)

	case "category":
		return Categories(body, path, method, user, idn, request)

	case "order":
		return Orders(body, path, method, user, idn, request)
	}

	fmt.Println("The request couldn't be processing.")

	return 400, METHOD_INVALID

}

func authorize(path string, method string, headers map[string]string) (bool, int, string) {

	// For this options the API does not have to validate the token.
	if (path == "product" && method == "GET") ||
		(path == "category" && method == "GET") {

		return true, 200, ""
	}

	// Get the token from headers.
	// The token to use must be the access_token in order to the validation process can get the Username attribute.
	token := headers["authorization"]
	if len(token) == 0 {
		return false, 401, "The token is mandatory to authenticate the requests."
	}

	// Validate token.
	// If the token is validated correctly, the msg variable will contain Username.
	ok, err, msg := auth.Validate(token)

	if !ok {
		if err != nil {
			errorMsg := err.Error()

			fmt.Println("Error", errorMsg)

			return false, 401, errorMsg

		} else {
			fmt.Println("Token error: ", msg)

		}
	}

	fmt.Println("Token is Ok.")

	return true, 200, msg

}

func Users(body string, path string, method string, user string, id string, request events.APIGatewayV2HTTPRequest) (int, string) {

	return 400, METHOD_INVALID
}

func Products(body string, path string, method string, user string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {

	return 400, METHOD_INVALID
}

func Categories(body string, path string, method string, user string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {

	switch method {
	case "POST":
		return routers.InsertCategory(body, user)
	case "PUT":
		return routers.UpdateCategory(body, user, id)
	}

	return 400, METHOD_INVALID
}

func Stock(body string, path string, method string, user string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {

	return 400, METHOD_INVALID
}

func Address(body string, path string, method string, user string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {

	return 400, METHOD_INVALID
}

func Orders(body string, path string, method string, user string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {

	return 400, METHOD_INVALID
}
