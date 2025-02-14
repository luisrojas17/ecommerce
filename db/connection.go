package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/luisrojas17/ecommerce/aws"
	"github.com/luisrojas17/ecommerce/models"
)

// It is defined this variable like pointer to set available this variable
// for all application in a better way.
var Connection *sql.DB

// This function gets a database connection.
func Connect() error {

	// To get the data connection to connect to database
	secretModel, err := aws.GetSecret(os.Getenv("SecretName"))

	if err != nil {
		fmt.Println("It was error to connect to database. It was not possible to get and parse the secret name.",
			err.Error())
		return err
	}

	Connection, err = sql.Open("mysql", buildConnectionString(secretModel))

	if err != nil {
		fmt.Println("It was an error when the application tried to connect to database.", err.Error())
		return err
	}

	err = Connection.Ping()
	if err != nil {
		fmt.Println("It was not possible to ping to database.", err.Error())
		return err
	}

	fmt.Println("Connection to database was successfull.")

	return nil
}

// This function wraps functionality to close database connection.
func Close() {
	fmt.Println("Closing connection...")

	Connection.Close()

	fmt.Println("Connection was closed.")
}

// This function wraps functionality to execute .
func Execute(statement string) error {

	fmt.Println(statement)

	_, err := Connection.Exec(statement)

	if err != nil {
		return err
	}

	return nil

}

// This function build the string connection.
// The string connection is built based on data connection which
// are wrapping into Secret structure.
func buildConnectionString(secret models.Secret) string {
	var user, pass, dbEndpoint, dbName string

	user = secret.Username
	pass = secret.Password
	dbEndpoint = secret.Host
	dbName = secret.DbName

	//connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?allowCleartextPassword=true", user, pass, dbEndpoint, dbName)

	// To specify to the driver to scan DATE and DATETIME automatically to time.Time
	// See: https://github.com/go-sql-driver/mysql#timetime-support
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", user, pass, dbEndpoint, dbName)

	fmt.Println("\nConnection string is: ", connectionString)

	return connectionString
}
