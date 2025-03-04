package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	"github.com/luisrojas17/ecommerce/models"
)

// Create an address in database. The address to create is related to the authenticated user.
func CreateAddress(address models.Address, userId string) (int64, error) {

	fmt.Println("Starting to insert address for user [" + userId + "] in database...")

	err := Connect()

	if err != nil {
		return 0, err
	}

	defer Close()

	statement := "INSERT INTO addresses (Add_UserId, Add_Address, Add_City, Add_State, Add_PostalCode, Add_Phone, Add_Title, Add_Name)"
	statement += " VALUES ('" + userId + "', '" + address.Street + "', '" + address.City + "', '" + address.State + "', '"
	statement += address.PostalCode + "', '" + address.Phone + "', '" + address.Title + "', '" + address.UserName + "')"

	var result sql.Result
	result, err = Connection.Exec(statement)
	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}

	idCreated, err2 := result.LastInsertId()
	if err2 != nil {
		return 0, err2
	}

	fmt.Println("Address was created sucessfully for user id: [" + userId + "].\nThe id created was: [" + strconv.Itoa(int(idCreated)) + "].")

	return idCreated, err2

}

// Update the address in database. The address to aupdate must be related to the authenticated user.
func UpdateAddress(address models.Address) error {

	fmt.Println("Starting to update address in database...")

	err := Connect()

	if err != nil {
		return err
	}

	defer Close()

	statement := "UPDATE addresses SET "

	if address.Street != "" {
		statement += "Add_Address = '" + address.Street + "', "
	}

	if address.City != "" {
		statement += "Add_City = '" + address.City + "', "
	}

	if address.State != "" {
		statement += "Add_State = '" + address.State + "', "
	}

	if address.PostalCode != "" {
		statement += "Add_PostalCode = '" + address.PostalCode + "', "
	}

	if address.Phone != "" {
		statement += "Add_Phone = '" + address.Phone + "', "
	}

	if address.Title != "" {
		statement += "Add_Title = '" + address.Title + "', "
	}

	if address.UserName != "" {
		statement += "Add_Name = '" + address.UserName + "', "
	}

	// We delete the last coma character in the sentence.
	statement, _ = strings.CutSuffix(statement, ", ")

	statement += "WHERE Add_Id = " + strconv.Itoa(address.Id)

	fmt.Println("Statement to execute: " + statement)

	_, err = Connection.Exec(statement)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Address id [" + strconv.Itoa(address.Id) + "] was updated sucessfully.")

	return nil

}

func ExistAddress(id int, userId string) (bool, error) {

	strIdAddress := strconv.Itoa(id)

	fmt.Println("Starting to search address id [" + strIdAddress + "] for user [" + userId + "] in database...")

	err := Connect()

	if err != nil {
		return false, err
	}

	defer Close()

	statement := "SELECT 1 FROM addresses WHERE Add_Id = " + strconv.Itoa(id) + " AND Add_UserId = '" + userId + "'"

	fmt.Println("Statement to execute: ", statement)

	rows, err := Connection.Query(statement)
	if err != nil {
		return false, err
	}

	var value string
	rows.Next()
	rows.Scan(&value)

	if value == "1" {
		fmt.Println("Address id [" + strIdAddress + "] was found for user [" + userId + "].")
		return true, nil
	}

	// If the address id and user id are not related to them.
	fmt.Println("Address id [" + strIdAddress + "] was not found for user [" + userId + "].")

	return false, nil
}

// This function delete an address in database. The address to be deleted has
// to match to parameter id provided.
func DeleteAddress(id int) (bool, error) {

	strId := strconv.Itoa(id)

	fmt.Println("Starting to delete address by id [" + strId + "] in database...")

	err := Connect()

	if err != nil {
		return false, err
	}

	defer Close()

	statement := "DELETE FROM addresses where Add_Id = " + strId

	var result sql.Result
	result, err = Connection.Exec(statement)
	if err != nil {
		fmt.Println(err.Error())
		return false, err
	}

	// We validated if there is any row affected.
	rowsAffected, err2 := result.RowsAffected()
	if err2 != nil {
		fmt.Println(err2.Error())
		return false, err2
	}

	if rowsAffected > 0 {
		fmt.Println("Address id [" + strId + "] was delete sucessfully. There was affected [" + strconv.Itoa(int(rowsAffected)) + "] rows.")
		return true, nil
	} else {
		fmt.Println("It was not possible to delete address id [" + strId + "]. There was affected [0] rows.")
		return false, nil
	}

}
