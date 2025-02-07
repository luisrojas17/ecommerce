package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"github.com/luisrojas17/ecommerce/models"
)

func CreateCategory(category models.Category) (int64, error) {

	fmt.Println("Starting to insert category in database...")

	err := Connect()

	if err != nil {
		return 0, err
	}

	defer Close()

	statement := "INSERT INTO category (Categ_Name, Categ_Path) VALUES('" + category.Name + "', '" + category.Path + "')"

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

	fmt.Println("Category was created sucessfully. The id created was: ", idCreated)

	return idCreated, err2
}
