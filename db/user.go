package db

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"github.com/luisrojas17/ecommerce/helper"
	"github.com/luisrojas17/ecommerce/models"
)

func CreateUser(user models.User) error {

	fmt.Println("Starting to insert user in database...")

	err := Connect()

	if err != nil {
		return err
	}

	defer Close()

	statement := "INSERT INTO users (User_Email, User_UUID, User_DateAdd) VALUES ('" + user.Email + "', '" + user.Uuid + "', '" + helper.GetDate() + "')"

	err = Execute(statement)
	if err != nil {
		return err
	}

	fmt.Println("User was inserted in database.")

	return nil
}

func IsAdmin(userId string) (bool, string) {
	fmt.Printf("Validating if user [%s] is admin.", userId)

	err := Connect()
	if err != nil {
		return false, err.Error()
	}

	defer Close()

	statement := "SELECT 1 FROM users WHERE User_UUID='" + userId + "' AND User_Status = 0"

	fmt.Println("Executing statement: ", statement)

	rows, err := Connection.Query(statement)

	if err != nil {
		return false, err.Error()
	}

	var value string

	// To start in first row.
	rows.Next()

	// Assign the value into variable
	rows.Scan(&value)

	fmt.Printf("The user [%s] is admin [%t].\n", userId, (value == "1"))

	if value == "1" {
		return true, "The user " + userId + " is administrator."
	} else {
		return false, "The user" + userId + " is not administrator."
	}

}
