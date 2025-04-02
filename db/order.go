package db

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/luisrojas17/ecommerce/models"
)

func CreateOrder(order models.Order) (int64, error) {

	fmt.Println("Starting to insert order in database...")

	err := Connect()

	if err != nil {
		return 0, err
	}

	defer Close()

	// To save order.
	statement := "INSERT INTO orders (Order_UserUUID, Order_Total, Order_AddId) VALUES('"

	statement += order.UserUuid + "'," + strconv.FormatFloat(order.Total, 'f', -1, 64) + "," + strconv.Itoa(order.AddressId) + ")"

	// result contains the last insert id and total rows affected and
	var result sql.Result
	result, err = Connection.Exec(statement)

	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}

	lastInsertId, err2 := result.LastInsertId()
	if err2 != nil {
		return 0, err2
	}

	// To save order's details.
	for _, details := range order.Details {
		statement = "INSERT INTO orders_detail (OD_OrderId, OD_ProdId, OD_Quantity, OD_Price) VALUES (" + strconv.Itoa(int(lastInsertId))

		statement += ", " + strconv.Itoa(details.ProductId) + "," + strconv.Itoa(details.Quantity) + "," + strconv.FormatFloat(details.Price, 'f', -1, 64) + ")"

		fmt.Println(statement)

		_, err = Connection.Exec(statement)
		if err != nil {
			fmt.Println(err.Error())
			return 0, err
		}
	}

	fmt.Println("Order [" + strconv.Itoa(int(lastInsertId)) + "] was inserted with all its details.")

	return lastInsertId, nil
}
