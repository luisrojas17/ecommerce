package db

import (
	"database/sql"
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

func UpdateUser(user models.User) error {

	fmt.Println("Starting to update user " + user.Uuid + " in database...")

	err := Connect()
	if err != nil {
		return err
	}

	defer Close()

	statement := "UPDATE users SET "

	semiColon := ""
	if len(user.FirstName) > 0 {
		semiColon = ", "
		statement += "User_FirstName = '" + user.FirstName + "'"
	}

	if len(user.LastName) > 0 {
		statement += semiColon + "User_LastName = '" + user.LastName + "'"
	}

	if len(semiColon) == 0 {
		semiColon = ", "
	}

	statement += semiColon + " User_DateUpg = '" + helper.GetDate() + "' WHERE User_UUID = '" + user.Uuid + "'"

	_, err = Connection.Exec(statement)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("User was updated sucessfully. The values updated were: [first name: " + user.FirstName + ", last name: " + user.LastName + "].")

	return nil
}

func GetUser(userId string) (models.User, error) {

	fmt.Println("Starting to get user by [id:" + userId + "] in database...")

	user := models.User{}

	err := Connect()
	if err != nil {
		return user, err
	}

	defer Close()

	statement := "SELECT * FROM users WHERE User_UUID = '" + userId + "'"

	var rows *sql.Rows
	rows, err = Connection.Query(statement)

	if err != nil {
		fmt.Println("It was an error to execute query: \n" + statement + ".\n" + err.Error())
		return user, err
	}

	defer rows.Close()

	rows.Next()

	var firstName sql.NullString
	var lastName sql.NullString
	var updateAt sql.NullTime

	rows.Scan(&user.Uuid, &user.Email, &firstName, &lastName, &user.Status, &user.CreatedAt, &updateAt)

	user.FirstName = firstName.String
	user.LastName = lastName.String
	user.UpdatedAt = updateAt.Time.String()

	fmt.Printf("\nThere were found user data for user id [%s].\n", userId)

	return user, nil
}

func ExistUser(userId string) (bool, string) {

	fmt.Println("Validating if user " + userId + " exist in database...")

	err := Connect()
	if err != nil {
		return false, err.Error()
	}

	defer Close()

	statement := "SELECT 1 FROM users WHERE User_UUID = '" + userId + "'"

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

	fmt.Printf("The user [%s] exist [%t].\n", userId, (value == "1"))

	if value == "1" {
		return true, "The user " + userId + " exist"
	} else {
		return false, "The user " + userId + " does not exist."
	}
}

func IsAdmin(userId string) (bool, string) {

	fmt.Printf("Validating if user [%s] is admin.", userId)

	err := Connect()
	if err != nil {
		return false, err.Error()
	}

	defer Close()

	statement := "SELECT 1 FROM users WHERE User_UUID = '" + userId + "' AND User_Status = 0"

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
