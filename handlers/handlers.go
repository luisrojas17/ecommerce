package handlers

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-lambda-go/events"

	"github.com/luisrojas17/ecommerce/auth"
)

// Returns the HTTP status code and message to describe the error.
func Handler(path string, method string, body string, headers map[string]string, request events.APIGatewayV2HTTPRequest) (int, string) {

	fmt.Println("Processing path: " + path + " for method: " + method)

	id := request.PathParameters["id"]

	// Convert the user id and product id. The product id only has to be number.
	idn, _ := strconv.Atoi(id)

	// Check if token is valid and it had not expired.
	ok, statusCode, user := authorize(path, method, body, headers)

	if !ok {
		return statusCode, user
	}

	switch path[0:4] {

	case "user":
		return Users(body, path, method, user, id, request)

	case "prod":
		return Products(body, path, method, user, idn, request)

	case "stoc":
		return Stock(body, path, method, user, idn, request)

	case "addr":
		return Address(body, path, method, user, idn, request)

	case "cate":
		return Categories(body, path, method, user, idn, request)

	case "orde":
		return Orders(body, path, method, user, idn, request)
	}

	return 400, "Method Invalid."

}

func authorize(path string, method string, body string, headers map[string]string) (bool, int, string) {

	// For this options the API does not have to validate the token.
	if (path == "product" && method == "GET") ||
		(path == "category" && method == "GET") {

		return true, 200, ""
	}

	// Get the token from headers.
	token := headers["authorization"]
	if len(token) == 0 {
		return false, 401, "The token is mandatory to authenticate the request.s"
	}

	// Validate token.
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

	fmt.Println("Token Ok.")

	return true, 200, msg

}

func Users(body string, path string, method string, user string, id string, request events.APIGatewayV2HTTPRequest) (int, string) {

	return 400, "Method Invalid"
}

func Products(body string, path string, method string, user string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {

	return 400, "Method Invalid"
}

func Categories(body string, path string, method string, user string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {

	return 400, "Method Invalid"
}

func Stock(body string, path string, method string, user string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {

	return 400, "Method Invalid"
}

func Address(body string, path string, method string, user string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {

	return 400, "Method Invalid"
}

func Orders(body string, path string, method string, user string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {

	return 400, "Method Invalid"
}
