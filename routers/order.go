package routers

import (
	"encoding/json"
	"net/http"
	"strconv"

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
