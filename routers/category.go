package routers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/luisrojas17/ecommerce/db"
	"github.com/luisrojas17/ecommerce/models"
)

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
		return 400, "It was an error to insert category [" + category.Name + "]." + err2.Error()
	}

	return 200, "{ categoryId: " + strconv.Itoa(int(result)) + " }"
}
