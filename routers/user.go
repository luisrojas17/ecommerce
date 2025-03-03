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

func UpdateUser(body string, userId string) (int, string) {

	var userToUpdate models.User

	err := json.Unmarshal([]byte(body), &userToUpdate)

	if err != nil {
		return http.StatusBadRequest, "It was an error to convert json user to model.\n" + err.Error()
	}

	if len(userToUpdate.FirstName) == 0 && len(userToUpdate.LastName) == 0 {
		return http.StatusPreconditionFailed, "You must specify User's first or last name."
	}

	existUser, _ := db.ExistUser(userId)

	if !existUser {
		return http.StatusNotFound, "The user " + userId + " does not exist."
	}

	userToUpdate.Uuid = userId

	err = db.UpdateUser(userToUpdate)
	if err != nil {
		return http.StatusBadRequest, "It was an error to update the user: " + userId + ".\n" + err.Error()
	}

	return http.StatusOK, "User Id: " + userId + " was updated successfully."
}

func GetUser(userId string) (int, string) {

	existUser, _ := db.ExistUser(userId)

	if !existUser {
		return http.StatusNotFound, "The user " + userId + " does not exist."
	}

	user, err := db.GetUser(userId)
	if err != nil {
		return http.StatusBadRequest, "It was an error to get user " + userId + ".\n" + err.Error()
	}

	jsonUser, err := json.Marshal(user)
	if err != nil {
		return http.StatusBadRequest, "It was an error to convert user to JSON format.\n" + err.Error()
	}

	return http.StatusOK, string(jsonUser)
}

func GetUsers(userId string, request events.APIGatewayV2HTTPRequest) (int, string) {

	var page, pageSize int
	var err error

	page, err = strconv.Atoi(request.QueryStringParameters["page"])
	if err != nil || page == 0 {
		page = 1
	}

	pageSize, err = strconv.Atoi(request.QueryStringParameters["pageSize"])
	if err != nil || pageSize > 25 {
		pageSize = 25
	}

	isAdmin, msg := db.IsAdmin(userId)
	if !isAdmin {
		fmt.Println("Only admin users can update product' stock.")
		return http.StatusPreconditionFailed, msg
	}

	var pageable models.PageableUsers
	pageable, err = db.GetUsers(page, pageSize)

	if err != nil {
		return http.StatusBadRequest, "It was an error to get users by [id:" + userId + "].\n" + err.Error()
	}

	var jsonPageable []byte
	jsonPageable, err = json.Marshal(pageable)
	if err != nil {
		return http.StatusBadRequest, "It was an error to convert users to JSON format.\n" + err.Error()
	}

	if pageable.TotalElements == 0 {
		return http.StatusNoContent, "There were not users related to user id: " + userId + "."
	}

	return http.StatusOK, string(jsonPageable)
}
