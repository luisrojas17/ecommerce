package db

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"github.com/luisrojas17/ecommerce/helper"
	"github.com/luisrojas17/ecommerce/models"
)

func Create(user models.User) error {

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
