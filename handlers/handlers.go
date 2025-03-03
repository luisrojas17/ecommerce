package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"

	"github.com/luisrojas17/ecommerce/auth"
	"github.com/luisrojas17/ecommerce/helper"
	"github.com/luisrojas17/ecommerce/routers"
)

// String constant to define a custom message to return
// when any API request can not be handled.
const METHOD_INVALID = "Method Invalid."

// This method handle all requests received by API Gagteway.
// Before to process the request, the method validate the request in order to know if
// the request will be processed or not. If the request contains a valid access_token
// the request will be passed accodding to path provided.
//
// Parameters:
//
//	path is an string which contains the path used in the URL. The possible pats can be:
//		category, product, stock, user.
//	method is an string which contains the method HTTP for which the request was sent it.
//	body is an string which contains the body sent by HTTP request. This is required for
//		create and update operations.
//	headers is an string which contains Authorization which store the access token
//		to authorize the request.
//	request is an events.APIGatewayV2HTTPRequest object which contains all data related to
//		request. For example: query parameters, path parameters
//
// Returns the HTTP status code and message to describe the error.
func Handler(path string, method string, body string, headers map[string]string, request events.APIGatewayV2HTTPRequest) (int, string) {

	fmt.Println("Processing request from path: " + path + " for method: " + method)

	id := request.PathParameters["id"]

	// Convert the user id and product id. The product id only has to be number.
	idn, _ := strconv.Atoi(id)

	// Check if token is valid and it had not expired.
	// Make sure to use the access_token in order to the validation can get the Username attribute.
	// The username is the id for AWS Cognito and this is the same value for database userUUID
	ok, statusCode, username := authorize(path, method, headers)

	if !ok {
		return statusCode, username
	}

	// To Do: process path in order to split from token "/" instead of use UrlPrefix variable.
	operation := strings.Split(path, "/")

	// We always get the first item from the array because that item contains the path
	// to identify the operation to do
	switch operation[0] {

	case "user":
		return Users(body, path, method, username, id, request)

	case "product":
		return Products(body, path, method, username, idn, request)

	case "stock":
		return Stock(body, path, method, username, idn, request)

	case "address":
		return Address(body, path, method, username, idn, request)

	case "category":
		return Categories(body, path, method, username, idn, request)

	case "order":
		return Orders(body, path, method, username, idn, request)

	default:
		fmt.Println("The request couldn't be processed.")

		return http.StatusBadRequest, METHOD_INVALID
	}

}

// This method validate the access_token in order to authorize the request received
// by the API Gateway. The access_token is contained in headers parameter like:
//
//	Authorization: access_token_value
//
// If the access_token is valid there will be returned:
//
//	boolean value to specify if the token is valid.
//	int value to specify HTTP code status. For example, if the request was process successful
//		and the access_token was valid the value to return will be 200. 401 for any other case.
//	string value which represents the username/user uuid if the token was valid. Error description
//		for any other case.
func authorize(path string, method string, headers map[string]string) (bool, int, string) {

	// For this options the API does not have to validate the token.
	if (path == "product" && method == "GET") ||
		(path == "category" && method == "GET") {

		return true, http.StatusOK, helper.EMPTY_STRING
	}

	// Get the token from headers.
	// The token to use must be the access_token in order to the validation process can get the Username attribute.
	token := headers["authorization"]
	if len(token) == 0 {
		return false, http.StatusUnauthorized, "The token is mandatory to authenticate the requests."
	}

	// Validate token.
	// If the token is validated correctly, the msg variable will contain Username.
	// Username is the id for AWS cognito. This is the same value for UserUUID in the database
	ok, err, msg := auth.Validate(token)

	if !ok {
		if err != nil {
			errorMsg := err.Error()

			fmt.Println("Error: ", errorMsg)

			return false, http.StatusUnauthorized, errorMsg

		} else {
			fmt.Println("Token error: ", msg)

		}
	}

	fmt.Println("Token is Ok.")

	return true, http.StatusOK, msg

}

// This method process all request related to user path.
// To consume this endpoint the request must contain an valid access_token.
func Users(body string, path string, method string, userId string, id string, request events.APIGatewayV2HTTPRequest) (int, string) {

	if path == "user/me" {
		switch method {
		case "PUT":
			//Check if it is possible to update status field.
			//HOST + /v1/ecommerce/user/me
			//	{
			//		"firstName": "Jose Luis",
			//		"lastName": "Rojas Gomez"
			//	}
			return routers.UpdateUser(body, userId)
		case "GET":
			//HOST + /v1/ecommerce/user/me
			return routers.GetUser(userId)
		default:
			fmt.Println("HTTP Method: [" + method + "] not found.")

			return http.StatusBadRequest, METHOD_INVALID
		}
	}

	if path == "user/all" {
		if method == "GET" {
			return routers.GetUsers(userId, request)
		}
	}

	return http.StatusBadRequest, METHOD_INVALID
}

// This API handled all CRUD operations related to products table.
func Products(body string, path string, method string, userId string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {

	switch method {
	case "POST":
		//HOST + /v1/ecommerce/product
		//	{
		//		"title": "Microwave",
		//		"description": "This is an appliance.",
		//		"price": 5000
		//	}
		return routers.CreateProduct(body, userId)
	case "PUT":
		// HOST + /v1/ecommerce/product/{id}
		//	{
		//		"categoryId": 8,
		//		"stock": 3
		//	}
		return routers.UpdateProduct(body, userId, id)
	case "DELETE":
		// HOST + /v1/ecommerce/product/{id}
		return routers.DeleteProduct(userId, id)
	case "GET":
		// HOST + /v1/ecommerce/product?pageSize=10&page=2&orderField=1&orderType=ASC&id=1&search=best&slug=&categoryId=7&categorySlug=Consoles
		return routers.GetProducts(request)
	default:
		fmt.Println("HTTP Method: [" + method + "] not found.")

		return http.StatusBadRequest, METHOD_INVALID
	}

}

// This API handled all CRUD operations related to category table.
func Categories(body string, path string, method string, userId string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {

	switch method {
	case "POST":
		// HOST + /v1/ecommerce/category
		//	{
		//		"name": "Video Games",
		//		"path": "Consoles"
		//	}
		return routers.CreateCategory(body, userId)
	case "PUT":
		// HOST + /v1/ecommerce/category/{id}
		//	{
		//		"name": "Instrumentos Musicales"
		//	}
		return routers.UpdateCategory(body, userId, id)
	case "DELETE":
		// HOST + /v1/ecommerce/category/{id}
		return routers.DeleteCategory(userId, id)
	case "GET":
		// HOST + /v1/ecommerce/category?slug=Television&categId=7
		return routers.GetCategories(request)
	default:
		fmt.Println("HTTP Method: [" + method + "] not found.")

		return http.StatusBadRequest, METHOD_INVALID
	}

}

// This API receive product id to update stock field in products table.
// Host + /v1/ecommerce/stock/{id}
func Stock(body string, path string, method string, userId string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {

	return routers.UpdateStock(body, userId, id)

}

func Address(body string, path string, method string, userId string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {

	switch method {
	case "POST":
		return routers.CreateAddress(body, userId)
	case "PUT":
		return routers.UpdateAddress(body, userId, id)
	default:
		fmt.Println("HTTP Method: [" + method + "] not found.")

		return http.StatusBadRequest, METHOD_INVALID
	}

}

func Orders(body string, path string, method string, userId string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {

	return http.StatusBadRequest, METHOD_INVALID
}
