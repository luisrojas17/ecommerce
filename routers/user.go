package routers

import (
	"encoding/json"

	"github.com/luisrojas17/ecommerce/db"
	"github.com/luisrojas17/ecommerce/models"
)

func UpdateUser(body string, userId string) (int, string) {

	var userToUpdate models.User

	err := json.Unmarshal([]byte(body), &userToUpdate)

	if err != nil {
		return 400, "It was an error to convert json user to model.\n" + err.Error()
	}

	if len(userToUpdate.FirstName) == 0 && len(userToUpdate.LastName) == 0 {
		return 400, "You must specify User's first or last name."
	}

	existUser, _ := db.ExistUser(userId)

	if !existUser {
		return 400, "The user " + userId + " does not exist."
	}

	userToUpdate.Uuid = userId

	err = db.UpdateUser(userToUpdate)
	if err != nil {
		return 400, "It was an error to update the user: " + userId + ".\n" + err.Error()
	}

	return 200, "User Id: " + userId + " was updated successfully."
}

func GetUser(userId string) (int, string) {

	existUser, _ := db.ExistUser(userId)

	if !existUser {
		return 400, "The user " + userId + " does not exist."
	}

	user, err := db.GetUser(userId)
	if err != nil {
		return 400, "It was an error to get user " + userId + ".\n" + err.Error()
	}

	jsonUser, err := json.Marshal(user)
	if err != nil {
		return 400, "It was an error to convert user to JSON format.\n" + err.Error()
	}

	return 200, string(jsonUser)
}
