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
		return 400, "It was an error to convert json user to model." + err.Error()
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
		return 400, "It was an error to update the user: " + userId + ". " + err.Error()
	}

	return 200, "User Id: " + userId + " was updated successfully."
}
