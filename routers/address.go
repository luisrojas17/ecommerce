package routers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/luisrojas17/ecommerce/db"
	"github.com/luisrojas17/ecommerce/models"
)

// Create an address in database. The address to aupdate must be related to the authenticated user.
func CreateAddress(body string, userId string) (int, string) {

	fmt.Println("Starting to create address for user id [" + userId + "].")

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
	return http.StatusOK, "{\"id\": " + strconv.Itoa(int(result)) + "}"

}

// Update the address in database. The address to aupdate must be related to the authenticated user.
func UpdateAddress(body string, userId string, idAddress int) (int, string) {

	strIdAddress := strconv.Itoa(idAddress)

	fmt.Println("Starting to update address id [" + strIdAddress + "] for user id [" + userId + "].")

	var address models.Address

	err := json.Unmarshal([]byte(body), &address)
	if err != nil {
		return http.StatusBadRequest, "It was an error to convert address' model to JSON format. " + err.Error()
	}

	// It is validate the address id for authenticate user id. This validation make sure that
	// the user id authenticated can only update his address.
	address.Id = idAddress
	var exist bool
	exist, err = db.ExistAddress(address.Id, userId)

	if !exist {
		if err != nil {
			return http.StatusBadRequest, "It was an error to search address id [" + strIdAddress + "] for user id [" + userId + "]. " + err.Error()
		}

		return http.StatusNotFound, "The address id [" + strIdAddress + "] was not found for user id [" + userId + "]."

	}

	err = db.UpdateAddress(address)
	if err != nil {
		return http.StatusBadRequest, "It was an error to update the address id [" + strIdAddress + "] for user id [" + userId + "]. " + err.Error()
	}

	return http.StatusOK, " The address id [" + strIdAddress + "] was updated for user id [" + userId + "]."

}

// Update the address in database. The address to aupdate must be related to the authenticated user.
func DeleteAddress(userId string, id int) (int, string) {

	if id == 0 {
		return http.StatusPreconditionFailed, "To delete any product you must specify the id since it is mandatory and it must be greater than 0."
	}

	// It is validate the address id for authenticate user id. This validation make sure that
	// the user id authenticated can only update his address.
	strIdAddress := strconv.Itoa(id)

	exist, err := db.ExistAddress(id, userId)

	if !exist {
		if err != nil {
			return http.StatusBadRequest, "It was an error to search address id [" + strIdAddress + "] for user id [" + userId + "]. " + err.Error()
		}

		return http.StatusNotFound, "The address id [" + strIdAddress + "] was not found for user id [" + userId + "]."

	}

	isDeleted, err := db.DeleteAddress(id)
	if err != nil {
		return http.StatusBadRequest, "It was an error to delete address [id: " + strIdAddress + "].\n" + err.Error()
	}

	if !isDeleted {
		return http.StatusBadRequest, "Address Id: " + strIdAddress + " does not exist. So, it could not be deleted."

	}

	return http.StatusOK, "Address Id: " + strIdAddress + " was deleted successfully."

}

func GetAddresses(userId string) (int, string) {

	addresses, err := db.GetAddresses(userId)
	if err != nil {
		return http.StatusBadRequest, "It was an error to get addresses list for user id [" + userId + "]. " + err.Error()
	}

	var jsonAddresses []byte
	jsonAddresses, err = json.Marshal(addresses)
	if err != nil {
		return http.StatusBadRequest, "It was an error to convert addresses list to JSON format.\n" + err.Error()
	}

	return http.StatusOK, string(jsonAddresses)
}
