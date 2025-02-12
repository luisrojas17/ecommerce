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

// Create a category in database.
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

// Update a category in database.
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

// This function delete a category in database. The product to be deleted has
// to match to parameter id provided.
func DeleteCategory(id int) (bool, error) {

	nId := strconv.Itoa(id)

	fmt.Println("Starting to delete category by id [" + nId + "] in database...")

	err := Connect()

	if err != nil {
		return false, err
	}

	defer Close()

	statement := "DELETE FROM category where Categ_Id = " + nId

	var result sql.Result
	result, err = Connection.Exec(statement)
	if err != nil {
		fmt.Println(err.Error())
		return false, err
	}

	rowsAffected, err2 := result.RowsAffected()
	if err2 != nil {
		fmt.Println(err2.Error())
		return false, err2
	}

	if rowsAffected > 0 {
		fmt.Println("Category id [" + nId + "] was delete sucessfully. There was affected [" + strconv.Itoa(int(rowsAffected)) + "] rows.")
		return true, nil
	} else {
		fmt.Println("It was not possible to delete category id [" + nId + "]. There was affected [0] rows.")
		return false, nil
	}

}

// Get a category in database. Initially the function return a category
// according to category id or slug provided like parameter. But If the id
// and slug parameters are not provided the function return all categories
// in the database.
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
