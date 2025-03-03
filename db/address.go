package db

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/go-sql-driver/mysql"

	"github.com/luisrojas17/ecommerce/models"
)

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
