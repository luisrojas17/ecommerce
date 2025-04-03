package db

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/luisrojas17/ecommerce/models"
)

func CreateOrder(order models.Order) (int64, error) {

	fmt.Println("Starting to insert order in database...")

	// To connect to database
	err := Connect()

	if err != nil {
		return 0, err
	}

	defer Close()

	// To start the transaction
	tx, err := Connection.BeginTx(context.Background(), nil)
	if err != nil {
		return 0, err
	}

	defer func() {
		// It does not matter if commit was completed.
		// If the commit was completed the Rollback will thrown an error but it is ignored.
		// On the other hand, in any case the Rollback will be completed
		_ = tx.Rollback()
	}()

	// To save order.
	statement := "INSERT INTO orders (Order_UserUUID, Order_Total, Order_AddId) VALUES('"

	statement += order.UserUuid + "'," + strconv.FormatFloat(order.Total, 'f', -1, 64) + "," + strconv.Itoa(order.AddressId) + ")"

	fmt.Println(statement)

	// result contains the last insert id and total rows affected and
	var result sql.Result
	//result, err = Connection.Exec(statement)
	result, err = tx.ExecContext(context.Background(), statement)

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

		//_, err = Connection.Exec(statement)
		_, err = tx.ExecContext(context.Background(), statement)
		if err != nil {
			fmt.Println(err.Error())
			return 0, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("It was an error to commit. %w", err)
	}

	fmt.Println("Order [" + strconv.Itoa(int(lastInsertId)) + "] was inserted with all its details.")

	return lastInsertId, nil
}

func GetOrders(userId string, dateFrom string, dateTo string, orderId int, page int, pageSize int) ([]models.Order, error) {

	fmt.Println("Starting to get order in database...")

	var orders []models.Order

	statement := "SELECT Order_Id, Order_UserUUID, Order_AddId, Order_Date, Order_Total FROM orders "

	// Fix this query because any user could get any orderId even if that order is not related to the user.
	if orderId > 0 {
		statement += " WHERE Order_Id = " + strconv.Itoa(orderId)
	} else {

		if page == 0 {
			page = 1
		}

		// To fix this assignment. If pagesize is not defined the value can be getting from environemnt variable.
		if pageSize == 0 {
			pageSize = 10
		}

		// To handle in the better way time in dates.
		// I mean, from the beginning of the day to the end of the day
		if len(dateFrom) == 10 {
			dateFrom += " 00:00:00"
		}

		if len(dateTo) == 10 {
			dateTo += " 23:59:59"
		}

		var where string

		// If we have valid dates, we will add them to where clause
		if len(dateFrom) > 0 && len(dateTo) > 0 {
			where += " WHERE Order_Date BETWEEN '" + dateFrom + "' AND '" + dateTo + "'"
		}

		// The user will always be there because the userId is gotten from token.
		// To make sure that the orders returned are related to the current user.
		var filterUserId string = " Order_UserUUID = '" + userId + "'"

		// We add the filter userId
		if len(where) > 0 {
			where += " AND " + filterUserId
		} else {
			where += " WHERE " + filterUserId
		}

		// Limit will be equals pageSize and define maximum records to get.
		limitAndOffset := " LIMIT " + strconv.Itoa(pageSize)

		// On the other hand, offset will be the record number where the query will start to get records.
		var offset int = (page * pageSize) - pageSize

		if offset > 0 {
			limitAndOffset += " OFFSET " + strconv.Itoa(offset)
		}

		statement += where + limitAndOffset

	}

	fmt.Println(statement)

	err := Connect()

	if err != nil {
		return orders, err
	}

	defer Close()

	var rows *sql.Rows
	rows, err = Connection.Query(statement)
	if err != nil {
		return orders, err
	}

	defer rows.Close()

	for rows.Next() {
		var order models.Order

		// To defined as null values becausae all of these variables
		// are optional values in the database schema.
		var total sql.NullFloat64

		// To set mandatory data for each order's property through a pointer
		err := rows.Scan(&order.Id, &order.UserUuid, &order.AddressId, &order.Date, &total)
		if err != nil {
			return orders, err
		}

		// To handle null values in temp variables
		order.Total = total.Float64

		// to get orders_detail
		var rowsDetails *sql.Rows
		statement := "SELECT OD_Id, OD_ProdId, OD_Quantity, OD_Price FROM orders_detail WHERE OD_OrderID = " + strconv.Itoa(order.Id)

		rowsDetails, err = Connection.Query(statement)
		if err != nil {
			return orders, err
		}

		for rowsDetails.Next() {

			// These variables do not require to handle null values because all of these
			// are primary and foreing key in the database schema
			var detailId int64
			var productId int64

			// On the other hand, next variables have to be defined as null values becausae all of these
			// are optional values in the database schema.
			var quantity sql.NullInt64
			var price sql.NullFloat64

			err = rowsDetails.Scan(&detailId, &productId, &quantity, &price)
			if err != nil {
				return orders, err
			}

			var details models.OrderDetails
			details.Id = int(detailId)
			details.ProductId = int(productId)
			details.Quantity = int(quantity.Int64)
			details.Price = price.Float64

			// To add each ordder's detail to the order gotten from database
			order.Details = append(order.Details, details)

		}

		// To add new item to array/slice orders
		orders = append(orders, order)

		rowsDetails.Close()
	}

	fmt.Println("There were found [" + strconv.Itoa(len(orders)) + "] orders.")

	return orders, nil

}
