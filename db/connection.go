package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/luisrojas17/ecommerce/aws"
	"github.com/luisrojas17/ecommerce/models"
)

// It is defined this variable like pointer to set available this variable for all application
// in a better way.
var Connection *sql.DB

func Connect() error {
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

func Close() {
	fmt.Println("Closing connection...")

	Connection.Close()

	fmt.Println("Connection was closed.")
}

func Execute(statement string) error {

	fmt.Println(statement)

	_, err := Connection.Exec(statement)

	if err != nil {
		return err
	}

	return nil

}

/*func Execute(statement string) (*sql.Rows, error) {

	fmt.Println(statement)

	rows, err := Connection.Query(statement)

	if err != nil {
		return nil, err
	}

	return rows, nil

}*/

func buildConnectionString(secret models.Secret) string {
	var user, pass, dbEndpoint, dbName string

	user = secret.Username
	pass = secret.Password
	dbEndpoint = secret.Host
	dbName = "ecommerce_db"

	//connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?allowCleartextPassword=true", user, pass, dbEndpoint, dbName)
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, pass, dbEndpoint, dbName)

	fmt.Println("Connection string is: ", connectionString)

	return connectionString
}
