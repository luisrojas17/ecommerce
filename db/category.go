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

func DeleteCategory(id int) error {

	nId := strconv.Itoa(id)

	fmt.Println("Starting to delete category by id [" + nId + "] in database...")

	err := Connect()

	if err != nil {
		return err
	}

	defer Close()

	statement := "DELETE FROM category where Categ_Id = " + nId

	_, err = Connection.Exec(statement)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Category id [" + nId + "] was delete sucessfully.")

	return nil
}

func GetCategories(id int, slug string) ([]models.Category, error) {

	nId := strconv.Itoa(id)

	fmt.Println("Starting to get categories by [id:" + nId + " and path: " + slug + "] in database...")

	// Define categies arrays.
	var categories []models.Category

	err := Connect()

	if err != nil {
		return categories, err
	}

	defer Close()

	statement := "SELECT * FROM category "

	if id > 0 {
		statement += "WHERE Categ_Id = " + strconv.Itoa(id)
	} else {
		if len(slug) > 0 {
			statement += "WHERE Categ_Path LIKE '%" + slug + "%'"
		}
	}

	fmt.Println("Statement to execute: " + statement)

	var rows *sql.Rows
	rows, err = Connection.Query(statement)

	if err != nil {
		fmt.Println("It was an error to execute query: \n" + statement + ".\n" + err.Error())
		return categories, err
	}

	for rows.Next() {
		var category models.Category

		// To handled null values
		var id sql.NullInt32
		var name sql.NullString
		var path sql.NullString

		err := rows.Scan(&id, &name, &path)
		if err != nil {
			fmt.Println("It was an error to iterate over the categories.\n" + err.Error())
			return categories, err
		}

		// To handled null values
		category.Id = int(id.Int32)
		category.Name = name.String
		category.Path = path.String

		// We add new item to array/slice categories
		categories = append(categories, category)
	}

	fmt.Printf("\nThere were found [%d] categories.\n", len(categories))

	return categories, nil

}
