package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	"github.com/luisrojas17/ecommerce/helper"
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

func UpdateCategory(category models.Category) error {

	fmt.Println("Starting to update category in database...")

	err := Connect()

	if err != nil {
		return err
	}

	defer Close()

	statement := "UPDATE category SET"

	if len(category.Name) > 0 {
		statement += " Categ_Name = '" + helper.RemoveQuotationMarks(category.Name) + "'"
	}

	if len(category.Path) > 0 {
		if !strings.HasSuffix(statement, "SET") {
			statement += ", "
		}

		statement += "Categ_Path = '" + helper.RemoveQuotationMarks(category.Path) + "'"
	}

	statement += " WHERE Categ_Id = '" + strconv.Itoa(category.Id) + "'"

	_, err = Connection.Exec(statement)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Category was updated sucessfully. The values updated were: [Id: " + strconv.Itoa(category.Id) + ", name: " + category.Name + ", path: " + category.Path + "].")

	return nil
}
