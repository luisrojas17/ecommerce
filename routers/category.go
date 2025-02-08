package routers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/luisrojas17/ecommerce/db"
	"github.com/luisrojas17/ecommerce/models"
)

// All functions defined in this file can only invoke by administrator users.

func InsertCategory(body string, User string) (int, string) {

	var category models.Category

	err := json.Unmarshal([]byte(body), &category)
	if err != nil {
		return 400, "It was an error to get category." + err.Error()
	}

	if len(category.Name) == 0 {
		return 400, "You must specify category's name."
	}

	if len(category.Path) == 0 {
		return 400, "You must specify category's path."
	}

	isAdmin, msg := db.IsAdmin(User)
	if !isAdmin {
		fmt.Println("Only admin users can create new catagories.")
		return 400, msg
	}

	result, err2 := db.CreateCategory(category)
	if err2 != nil {
		return 400, "It was an error to insert category [" + category.Name + "].\n" + err2.Error()
	}

	return 200, "{ categoryId: " + strconv.Itoa(int(result)) + " }"
}

func UpdateCategory(body string, User string, id int) (int, string) {

	var category models.Category

	// Convert json body to struct category
	err := json.Unmarshal([]byte(body), &category)

	if err != nil {
		return 400, "It was an error to get category." + err.Error()
	}

	if len(category.Name) == 0 && len(category.Path) == 0 {
		return 400, "You must specify category's name and category's path."
	}

	isAdmin, msg := db.IsAdmin(User)
	if !isAdmin {
		fmt.Println("Only admin users can create new catagories.")
		return 400, msg
	}

	category.Id = id

	err2 := db.UpdateCategory(category)
	if err2 != nil {
		return 400, "It was an error to update category [id: " + strconv.Itoa(id) + ", name: " + category.Name + "].\n" + err2.Error()
	}

	return 200, "CategoryId: " + strconv.Itoa(id) + " was updated successfully."
}

func DeleteCategory(User string, id int) (int, string) {

	if id == 0 {
		return 412, "To delete any category you must specify the id since it is mandatory and it must be greater than 0."
	}

	isAdmin, msg := db.IsAdmin(User)
	if !isAdmin {
		fmt.Println("Only admin users can create new catagories.")
		return 400, msg
	}

	err := db.DeleteCategory(id)
	if err != nil {
		return 400, "It was an error to delete category id: " + strconv.Itoa(id) + ".\n" + err.Error()
	}

	return 200, "CategoryId: " + strconv.Itoa(id) + " was deleted successfully."
}
