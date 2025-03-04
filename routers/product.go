package routers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"

	"github.com/luisrojas17/ecommerce/db"
	"github.com/luisrojas17/ecommerce/helper"
	"github.com/luisrojas17/ecommerce/models"
)

func CreateProduct(body string, userId string) (int, string) {

	var product models.Product

	err := json.Unmarshal([]byte(body), &product)
	if err != nil {
		return http.StatusBadRequest, "It was an error to convert JSON product to model. " + err.Error()
	}

	// We only validate this field since this is mandatory field. All other fields are optionals.
	if len(product.Title) == 0 {
		return http.StatusBadRequest, "You must specify product's title."
	}

	isAdmin, msg := db.IsAdmin(userId)
	if !isAdmin {
		fmt.Println("Only admin users can create new products.")
		return http.StatusBadRequest, msg
	}

	result, err2 := db.CreateProduct(product)

	if err2 != nil {
		return http.StatusBadRequest, "It was an error to insert product [" + product.Title + "].\n" + err2.Error()
	}

	// We return the Id generated for product inserted.
	return http.StatusOK, "{ id: " + strconv.Itoa(int(result)) + " }"
}

func UpdateProduct(body string, userId string, id int) (int, string) {

	var product models.Product

	// Convert json body to struct product
	err := json.Unmarshal([]byte(body), &product)

	if err != nil {
		return http.StatusBadRequest, "It was an error to convert product' model to JSON format. " + err.Error()
	}

	isAdmin, msg := db.IsAdmin(userId)
	if !isAdmin {
		fmt.Println("Only admin users can update products.")
		return http.StatusBadRequest, msg
	}

	product.Id = id

	err2 := db.UpdateProduct(product)
	if err2 != nil {
		return http.StatusBadRequest, "It was an error to update product [id: " + strconv.Itoa(id) + "].\n" + err2.Error()
	}

	return http.StatusOK, "Product Id: " + strconv.Itoa(id) + " was updated successfully."
}

func DeleteProduct(userId string, id int) (int, string) {

	if id == 0 {
		return http.StatusPreconditionFailed, "To delete any product you must specify the id since it is mandatory and it must be greater than 0."
	}

	isAdmin, msg := db.IsAdmin(userId)
	if !isAdmin {
		fmt.Println("Only admin users can delete products.")
		return http.StatusBadRequest, msg
	}

	strId := strconv.Itoa(id)
	isDeleted, err := db.DeleteProduct(id)
	if err != nil {
		return http.StatusBadRequest, "It was an error to delete product [id: " + strId + "].\n" + err.Error()
	}

	if !isDeleted {
		return http.StatusBadRequest, "Product Id: " + strId + " does not exist. So, it could not be deleted."

	}

	return http.StatusOK, "Product Id: " + strId + " was deleted successfully."
}

// This function get all products in a paged way.
func GetProducts(request events.APIGatewayV2HTTPRequest) (int, string) {

	var product models.Product
	var page, pageSize int
	var orderType, orderField string

	queryPameters := request.QueryStringParameters

	fmt.Println("QueryParameters to process get products request: ", queryPameters)

	page, _ = strconv.Atoi(queryPameters["page"])
	pageSize, _ = strconv.Atoi(queryPameters["pageSize"])
	orderType = queryPameters["orderType"]   // DESC or ASC
	orderField = queryPameters["orderField"] // Id = 1, Title = 2, Description = 3, CreatedAt = 4,
	// Price = 5, Stock = 6, CategoryId = 7

	if !strings.Contains("123456", orderField) {
		orderField = helper.EMPTY_STRING
	}

	//var choice string
	//if len(queryPameters["id"]) > 0 {choice = "ID"}
	product.Id, _ = strconv.Atoi(queryPameters["id"])
	product.Search = queryPameters["search"]
	product.Path = queryPameters["slug"]
	product.CategoryId, _ = strconv.Atoi(queryPameters["categoryId"])
	product.CategoryPath = queryPameters["categorySlug"]

	var pageable models.Pageable
	var err error
	pageable, err = db.GetProducts(product, page, pageSize, orderType, orderField)

	if err != nil {
		return http.StatusBadRequest, "It was an error to get products.\n" + err.Error()
	}

	var jsonPageable []byte
	jsonPageable, err = json.Marshal(pageable)
	if err != nil {
		return http.StatusBadRequest, "It was an error to convert products to JSON format.\n" + err.Error()
	}

	return http.StatusOK, string(jsonPageable)

}

func UpdateStock(body string, userId string, id int) (int, string) {

	var product models.Product

	// Convert json body to struct product
	err := json.Unmarshal([]byte(body), &product)

	if err != nil {
		return http.StatusBadRequest, "It was an error to convert product stock's model to JSON format. " + err.Error()
	}

	isAdmin, msg := db.IsAdmin(userId)
	if !isAdmin {
		fmt.Println("Only admin users can update product' stock.")
		return http.StatusBadRequest, msg
	}

	product.Id = id

	err2 := db.UpdateStock(product)
	if err2 != nil {
		return http.StatusBadRequest, "It was an error to update product's stoke [id: " + strconv.Itoa(id) + "].\n" + err2.Error()
	}

	return http.StatusOK, "Product id stock: " + strconv.Itoa(id) + " was updated successfully."
}
