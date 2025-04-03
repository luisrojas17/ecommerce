package routers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/luisrojas17/ecommerce/db"
	"github.com/luisrojas17/ecommerce/models"
)

func CreateOrder(body string, userId string) (int, string) {

	var order models.Order

	err := json.Unmarshal([]byte(body), &order)

	if err != nil {
		return http.StatusBadRequest, "It was an error to convert JSON order to model. " + err.Error()
	}

	order.UserUuid = userId

	ok, message := validateOrder(order)

	if !ok {
		return http.StatusBadRequest, message
	}

	result, err2 := db.CreateOrder(order)

	if err2 != nil {
		return http.StatusBadRequest, "It was an error to create the order for user id [" + userId + "]." + err2.Error()
	}

	return http.StatusOK, "{\"id\": " + strconv.Itoa(int(result)) + "}"

}

func validateOrder(order models.Order) (bool, string) {

	if order.Total == 0 {
		return false, "You must set the total orders."
	}

	count := 0
	for _, orderDetails := range order.Details {

		if orderDetails.ProductId == 0 {
			return false, "You must set the product id."
		}

		if orderDetails.Quantity == 0 {
			return false, "You must set the quantity."
		}

		count++
	}

	if count == 0 {
		return false, "You must set order details."
	}

	return true, ""
}

func GetOrders(userId string, request events.APIGatewayV2HTTPRequest) (int, string) {

	var dateFrom, dateTo string
	var orderId int
	var page int
	var pageSize int

	queryPameters := request.QueryStringParameters

	fmt.Println("QueryParameters to process get orders request: ", queryPameters)

	// We validate page and pageSize in database layer to calculate limit and offset
	page, _ = strconv.Atoi(queryPameters["page"])
	pageSize, _ = strconv.Atoi(queryPameters["pageSize"])

	// It does not matter if there were any error to convert the queryparameters
	// since we can use the default values for each one of parameters
	orderId, _ = strconv.Atoi(queryPameters["id"])

	// Check if we need to validate dates because these values will be necessary for query.
	dateFrom = queryPameters["dateFrom"]
	dateTo = queryPameters["dateTo"]

	result, err := db.GetOrders(userId, dateFrom, dateTo, orderId, page, pageSize)
	if err != nil {
		return http.StatusBadGateway, "It was an error to get orders with next parameters: [dateFrom: " + dateFrom + ", dateTo: " + dateTo + "]. Error: " + err.Error()
	}

	if len(result) == 0 {
		return http.StatusNoContent, "There were not orders related to user id: " + userId + "."
	}

	// Convert model to JSON byte array
	Order, err2 := json.Marshal(result)
	if err2 != nil {
		return http.StatusBadGateway, "It was an error to convert Order model to JSON string. Error: " + err2.Error()
	}

	return http.StatusOK, string(Order)
}
