package routers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/luisrojas17/ecommerce/db"
	"github.com/luisrojas17/ecommerce/models"
)

func CreateAddress(body string, userId string) (int, string) {

	var address models.Address

	err := json.Unmarshal([]byte(body), &address)
	if err != nil {
		return http.StatusBadRequest, "It was an error to convert JSON address to model. " + err.Error()
	}

	// Starting with validations
	if address.Title == "" {
		return http.StatusPreconditionFailed, "You must specify address's title."
	}

	if address.Street == "" {
		return http.StatusPreconditionFailed, "You must specify address's street."
	}

	if address.City == "" {
		return http.StatusPreconditionFailed, "You must specify address' city."
	}

	if address.PostalCode == "" {
		return http.StatusPreconditionFailed, "You must specify address's postal code."
	}

	if address.Phone == "" {
		return http.StatusBadRequest, "You must specify address's phone."
	}

	if address.UserName == "" {
		return http.StatusPreconditionFailed, "You must specify address's user name."
	}

	var result int64
	result, err = db.CreateAddress(address, userId)

	if err != nil {
		return http.StatusBadRequest, "It was an error to create the address for user id [" + userId + "]. " + err.Error()
	}

	// Check if it is better using status code http.StatusCreated
	return http.StatusOK, "{id: " + strconv.Itoa(int(result)) + "}"

}
