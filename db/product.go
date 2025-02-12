package db

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/luisrojas17/ecommerce/helper"
	"github.com/luisrojas17/ecommerce/models"
)

func CreateProduct(product models.Product) (int64, error) {

	fmt.Println("Starting to insert product in database...")

	err := Connect()

	if err != nil {
		return 0, err
	}

	defer Close()

	statement := "INSERT INTO products (Prod_Title "

	if len(product.Description) > 0 {
		statement += ", Prod_Description"
	}

	if product.Price > 0 {
		statement += ", Prod_Price"
	}

	if product.CategoryId > 0 {
		statement += ", Prod_CategoryId"
	}

	if product.Stock > 0 {
		statement += ", Prod_Stock"
	}

	if len(product.Path) > 0 {
		statement += ", Prod_Path"
	}

	statement += ") VALUES ('" + helper.RemoveQuotationMarks(product.Title) + "'"

	if len(product.Description) > 0 {
		statement += ",'" + helper.RemoveQuotationMarks(product.Description) + "'"
	}

	if product.Price > 0 {
		statement += ", " + strconv.FormatFloat(product.Price, 'e', -1, 64)
	}

	if product.CategoryId > 0 {
		statement += ", " + strconv.Itoa(product.CategoryId)
	}

	if product.Stock > 0 {
		statement += ", " + strconv.Itoa(product.Stock)
	}

	if len(product.Path) > 0 {
		statement += ", '" + helper.RemoveQuotationMarks(product.Path) + "'"
	}

	statement += ")"

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

	fmt.Println("Product was created sucessfully. The id created was: ", idCreated)

	return idCreated, err2
}

func UpdateProduct(product models.Product) error {

	fmt.Println("Starting to update category in database...")

	err := Connect()

	if err != nil {
		return err
	}

	defer Close()

	statement := "UPDATE products SET "

	statement = helper.BuildStatement(statement, "Prod_Title", "S", 0, 0, product.Title)
	statement = helper.BuildStatement(statement, "Prod_Description", "S", 0, 0, product.Description)
	statement = helper.BuildStatement(statement, "Prod_Price", "F", 0, product.Price, helper.EMPTY_STRING)
	statement = helper.BuildStatement(statement, "Prod_CategoryId", "N", product.CategoryId, 0, helper.EMPTY_STRING)
	statement = helper.BuildStatement(statement, "Prod_Stock", "N", product.Stock, 0, helper.EMPTY_STRING)
	statement = helper.BuildStatement(statement, "Prod_Path", "S", 0, 0, product.Path)

	statement += " WHERE Prod_Id = " + strconv.Itoa(product.Id)

	fmt.Println("Statement to execute: " + statement)

	_, err = Connection.Exec(statement)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Product id [" + strconv.Itoa(product.Id) + "] was updated sucessfully.")

	return nil
}
