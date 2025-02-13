package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

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

// This function delete a product in database. The product to be deleted has
// to match to parameter id provided.
func DeleteProduct(id int) (bool, error) {

	nId := strconv.Itoa(id)

	fmt.Println("Starting to delete product by id [" + nId + "] in database...")

	err := Connect()

	if err != nil {
		return false, err
	}

	defer Close()

	statement := "DELETE FROM products WHERE Prod_Id = " + nId

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
		fmt.Println("Product id [" + nId + "] was delete sucessfully. There was affected [" + strconv.Itoa(int(rowsAffected)) + "] rows.")
		return true, nil
	} else {
		fmt.Println("It was not possible to delete product id [" + nId + "]. There was affected [0] rows.")
		return false, nil
	}

}

func GetProducts(product models.Product, page int, pageSize int, orderType string, orderField string) (models.Pageable, error) {

	fmt.Printf("\nStarting to get products %#v from in database...", product)

	var pageable models.Pageable
	var products []models.Product

	err := Connect()

	if err != nil {
		return pageable, err
	}

	defer Close()

	var countStatement string
	var statement string
	var where string

	countStatement = "SELECT count(*) as records FROM products "

	statement = "SELECT * FROM products "

	// To prepare conditions for where clause.
	if product.Id > 0 {
		where = " WHERE Prod_Id = " + strconv.Itoa(product.Id)
	} else if product.Search != helper.EMPTY_STRING {
		where = " WHERE UCASE(CONCAT(Prod_Title, Prod_Description)) LIKE '%" + strings.ToUpper(product.Search) + "%' "
	} else if product.CategoryId > 0 {
		where = " WHERE Prod_CategoryId " + strconv.Itoa(product.CategoryId)
	} else if product.Path != helper.EMPTY_STRING {
		where = " WHERE UCASE(Prod_Path) LIKE '%" + strings.ToUpper(product.Path) + "%' "
	} else if product.CategoryPath != helper.EMPTY_STRING {
		join := " JOIN category ON Prod_CategoryId = Categ_Id AND Categ_Path LIKE '%" + strings.ToUpper(product.CategoryPath) + "%' "
		countStatement += join
		statement += join
	} else {
		fmt.Println("There was not any search criteria.")
	}

	countStatement += where

	fmt.Println("Count statement to execute: " + countStatement)

	// To execute count statement.
	var rows *sql.Rows
	rows, err = Connection.Query(countStatement)

	if err != nil {
		fmt.Println("It was an error to execute query: \n" + statement + ".\n" + err.Error())
		return pageable, err
	}

	defer rows.Close()

	// To get total records.
	rows.Next()
	var totalRecords sql.NullInt32
	err = rows.Scan(&totalRecords)

	if err != nil {
		fmt.Println("It was an error to get the total records for query: \n" + statement + ".\n" + err.Error())
		return pageable, err
	}

	// To assign total records to peagable struct.
	totalElments := int(totalRecords.Int32)
	pageable.TotalElements = totalElments

	fmt.Println("There were found [" + strconv.Itoa(totalElments) + "] products.")

	// If there were not any result the flow has to break to avoid continuing.
	if totalElments <= 0 {
		return pageable, nil
	}

	// To calculate limit and offset to page data.
	var limit string
	if page > 0 {
		if totalElments > pageSize {
			limit = " LIMIT " + strconv.Itoa(pageSize)
			if page > 1 {
				offset := pageSize * (page - 1)
				limit += " OFFSET " + strconv.Itoa(offset)
			}
		} else {
			limit = helper.EMPTY_STRING
		}
	}

	// To build order by clause.
	var orderBy string
	if len(orderField) > 0 {
		// orderField can contain someone this values:
		// Id = 1, Title = 2, Description = 3, CreatedAt = 4, Price = 5, Stock = 6, CategoryId = 7
		switch orderField {
		case "1":
			orderBy = " ORDER BY Prod_Id "
		case "2":
			orderBy = " ORDER BY Prod_Title "
		case "3":
			orderBy = " ORDER BY Prod_Description "
		case "4":
			orderBy = " ORDER BY Prod_CreatedAt "
		case "5":
			orderBy = " ORDER BY Prod_Price "
		case "6":
			orderBy = " ORDER BY Prod_Stock "
		case "7":
			orderBy = " ORDER BY Prod_CategoryId "
		default:
			orderBy = helper.EMPTY_STRING
		}

		if orderBy != helper.EMPTY_STRING && orderType == models.ASC || orderType == models.DESC {
			orderBy += " " + orderType
		}

	}

	// To build statement with where, orderBy and limit clauses.
	statement += where + orderBy + limit

	fmt.Println("Statement to execute: " + statement)

	// To execute statement to get all products in database.
	rows, err = Connection.Query(statement)

	if err != nil {
		fmt.Println("It was an error to execute query: \n" + statement + ".\n" + err.Error())
		return pageable, err
	}

	// To iterate over result.
	for rows.Next() {
		var product models.Product

		// To handled null values
		var id sql.NullInt32
		var title sql.NullString
		var description sql.NullString
		var createdAt sql.NullTime
		var updatedAt sql.NullTime
		var price sql.NullFloat64
		var path sql.NullString
		var categoryId sql.NullInt32
		var stock sql.NullInt32

		// To map each column to variable defined before
		err := rows.Scan(&id, &title, &description, &createdAt, &updatedAt, &price, &path, &categoryId, &stock)
		if err != nil {
			fmt.Println("It was an error to map values to columns for products.\n" + err.Error())
			return pageable, err
		}

		product.Id = int(id.Int32)
		product.Title = title.String
		product.Description = description.String
		product.CreatedAt = createdAt.Time.String()
		product.UpdatedAt = updatedAt.Time.String()
		product.Price = price.Float64
		product.Path = path.String
		product.CategoryId = int(categoryId.Int32)
		product.Stock = int(stock.Int32)

		// We add new item to array/slice categories
		products = append(products, product)
	}

	fmt.Printf("\nThere were found [%d] products.\n", len(products))

	// To add all products to struct pageable
	pageable.Content = products

	return pageable, nil
}
