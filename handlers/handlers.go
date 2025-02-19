package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"

	"github.com/luisrojas17/ecommerce/auth"
	"github.com/luisrojas17/ecommerce/routers"
)

// String constant to define a custom message to return
// when any API request can not be handled.
const METHOD_INVALID = "Method Invalid."

// Returns the HTTP status code and message to describe the error.
func Handler(path string, method string, body string, headers map[string]string, request events.APIGatewayV2HTTPRequest) (int, string) {

	fmt.Println("Processing request from path: " + path + " for method: " + method)

	id := request.PathParameters["id"]

	// Convert the user id and product id. The product id only has to be number.
	idn, _ := strconv.Atoi(id)

	// Check if token is valid and it had not expired.
	// Make sure to use the access_token in order to the validation can get the Username attribute.
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

	default:
		fmt.Println("The request couldn't be processing.")

		return 400, METHOD_INVALID
	}

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

// This API handled all CRUD operations related to products table.
func Products(body string, path string, method string, user string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {

	switch method {
	case "POST":
		//HOST + /v1/ecommerce/product
		//	{
		//		"title": "Microwave",
		//		"description": "This is an appliance.",
		//		"price": 5000
		//	}
		return routers.CreateProduct(body, user)
	case "PUT":
		// HOST + /v1/ecommerce/product/{id}
		//	{
		//		"categoryId": 8,
		//		"stock": 3
		//	}
		return routers.UpdateProduct(body, user, id)
	case "DELETE":
		// HOST + /v1/ecommerce/product/{id}
		return routers.DeleteProduct(user, id)
	case "GET":
		// HOST + /v1/ecommerce/product?pageSize=10&page=2&orderField=1&orderType=ASC&id=1&search=best&slug=&categoryId=7&categorySlug=Consoles
		return routers.GetProducts(request)
	default:
		fmt.Println("HTTP Method: [" + method + "] not found.")

		return 400, METHOD_INVALID
	}

}

// This API handled all CRUD operations related to category table.
func Categories(body string, path string, method string, user string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {

	switch method {
	case "POST":
		// HOST + /v1/ecommerce/category
		//	{
		//		"name": "Video Games",
		//		"path": "Consoles"
		//	}
		return routers.CreateCategory(body, user)
	case "PUT":
		// HOST + /v1/ecommerce/category/{id}
		//	{
		//		"name": "Instrumentos Musicales"
		//	}
		return routers.UpdateCategory(body, user, id)
	case "DELETE":
		// HOST + /v1/ecommerce/category/{id}
		return routers.DeleteCategory(user, id)
	case "GET":
		// HOST + /v1/ecommerce/category?slug=Television&categId=7
		return routers.GetCategories(request)
	default:
		fmt.Println("HTTP Method: [" + method + "] not found.")

		return 400, METHOD_INVALID
	}

}

// This API receive product id to update stock field in products table.
// Host + /v1/ecommerce/stock/{id}
func Stock(body string, path string, method string, user string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {

	return routers.UpdateStock(body, user, id)

}

func Address(body string, path string, method string, user string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {

	return 400, METHOD_INVALID
}

func Orders(body string, path string, method string, user string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {

	return 400, METHOD_INVALID
}
