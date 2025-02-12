package routers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/luisrojas17/ecommerce/db"
	"github.com/luisrojas17/ecommerce/models"
)

// All functions defined in this file can only invoke by administrator users.

func CreateCategory(body string, user string) (int, string) {

	var category models.Category

	err := json.Unmarshal([]byte(body), &category)
	if err != nil {
		return 400, "It was an error to convert json category to model.  " + err.Error()
	}

	if len(category.Name) == 0 {
		return 400, "You must specify category's name."
	}

	if len(category.Path) == 0 {
		return 400, "You must specify category's path."
	}

	isAdmin, msg := db.IsAdmin(user)
	if !isAdmin {
		fmt.Println("Only admin users can create new catagories.")
		return 400, msg
	}

	result, err2 := db.CreateCategory(category)
	if err2 != nil {
		return 400, "It was an error to insert category [" + category.Name + "].\n" + err2.Error()
	}

	return 200, "{ id: " + strconv.Itoa(int(result)) + " }"
}

func UpdateCategory(body string, user string, id int) (int, string) {

	var category models.Category

	// Convert json body to struct category
	err := json.Unmarshal([]byte(body), &category)

	if err != nil {
		return 400, "It was an error to get category. " + err.Error()
	}

	if len(category.Name) == 0 && len(category.Path) == 0 {
		return 400, "You must specify category's name and category's path."
	}

	isAdmin, msg := db.IsAdmin(user)
	if !isAdmin {
		fmt.Println("Only admin users can create new catagories.")
		return 400, msg
	}

	category.Id = id

	err2 := db.UpdateCategory(category)
	if err2 != nil {
		return 400, "It was an error to update category [id: " + strconv.Itoa(id) + ", name: " + category.Name + "].\n" + err2.Error()
	}

	return 200, "Category Id: " + strconv.Itoa(id) + " was updated successfully."
}

func DeleteCategory(user string, id int) (int, string) {

	if id == 0 {
		return 412, "To delete any category you must specify the id since it is mandatory and it must be greater than 0."
	}

	isAdmin, msg := db.IsAdmin(user)
	if !isAdmin {
		fmt.Println("Only admin users can create new catagories.")
		return 400, msg
	}

	isDeleted, err := db.DeleteCategory(id)
	if err != nil {
		return 400, "It was an error to delete category id: " + strconv.Itoa(id) + ".\n" + err.Error()
	}

	if !isDeleted {
		return 404, "Category Id: " + strconv.Itoa(id) + " does not exist. So, it could not be deleted."

	}

	return 200, "Category Id: " + strconv.Itoa(id) + " was deleted successfully."
}

func GetCategories(request events.APIGatewayV2HTTPRequest) (int, string) {

	var err error
	var categoryId int
	var strCategoryId string
	var slug string

	strCategoryId = request.QueryStringParameters["id"]

	if len(strCategoryId) > 0 {
		categoryId, err = strconv.Atoi(strCategoryId)

		if err != nil {
			return 412, "The category id must be a numeric value greater than 0. " + err.Error()
		}
	} else {
		if len(request.QueryStringParameters["slug"]) > 0 {
			slug = request.QueryStringParameters["slug"]
		}
	}

	categories, err := db.GetCategories(categoryId, slug)

	if err != nil {
		return 400, "It was an error to get category by [id:" + strCategoryId + " and Path:" + slug + "].\n" + err.Error()
	}

	if len(categories) == 0 {
		return 201, "There was not categories related to category id: " + strCategoryId + " or Path: " + slug
	}

	jsonCategories, err := json.Marshal(categories)
	if err != nil {
		return 400, "It was an error to convert categories to JSON format.\n" + err.Error()
	}

	return 200, string(jsonCategories)
}
