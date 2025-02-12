package routers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/luisrojas17/ecommerce/db"
	"github.com/luisrojas17/ecommerce/models"
)

func CreateProduct(body string, User string) (int, string) {

	var product models.Product

	err := json.Unmarshal([]byte(body), &product)
	if err != nil {
		return 400, "It was an error to convert json product to model. " + err.Error()
	}

	// We only validate this field since this is mandatory field. All other fields are optionals.
	if len(product.Title) == 0 {
		return 400, "You must specify product's title."
	}

	isAdmin, msg := db.IsAdmin(User)
	if !isAdmin {
		fmt.Println("Only admin users can create new catagories.")
		return 400, msg
	}

	result, err2 := db.CreateProduct(product)

	if err2 != nil {
		return 400, "It was an error to insert product [" + product.Title + "].\n" + err2.Error()
	}

	// We return the Id generated for product inserted.
	return 200, "{ id: " + strconv.Itoa(int(result)) + " }"
}

func UpdateProduct(body string, User string, id int) (int, string) {

	var product models.Product

	// Convert json body to struct category
	err := json.Unmarshal([]byte(body), &product)

	if err != nil {
		return 400, "It was an error to get product. " + err.Error()
	}

	isAdmin, msg := db.IsAdmin(User)
	if !isAdmin {
		fmt.Println("Only admin users can create new catagories.")
		return 400, msg
	}

	product.Id = id

	err2 := db.UpdateProduct(product)
	if err2 != nil {
		return 400, "It was an error to update product [id: " + strconv.Itoa(id) + "].\n" + err2.Error()
	}

	return 200, "Product Id: " + strconv.Itoa(id) + " was updated successfully."
}
