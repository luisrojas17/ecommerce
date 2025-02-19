package aws

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"

	"github.com/luisrojas17/ecommerce/models"
)

// This function gets the secret configurated in AWS Secret Manager.
// The secret contains all data related to database connection which
// will be wrapping into Secret Model.
// This is an example what this function will get when is invoked the
// AWS Secrets Manager API:
//
//	{
//		"username":"somevalue",
//		"password":"somevalue",
//		"engine":"mysql",
//		"host":"somevalue",
//		"port":somevalue,
//		"dbname":"somevalue",
//		"dbInstanceIdentifier":"somevalue"
//	}
//
// Param "secretName" is a name of the secretconfigure in AWS Secrets Manager.
// This value is defined like environment variable.
// See: Lambda configuration -> Environment variables
func GetSecret(secretName string) (models.Secret, error) {
	var secretModel models.Secret

	fmt.Printf("\nRequesting secretName [%s]...", secretName)

	svc := secretsmanager.NewFromConfig(Cfg)

	// Key is an json structure
	key, err := svc.GetSecretValue(Ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	})

	if err != nil {
		fmt.Println("It was an error to retrieve the secret name and parsing it. ", err.Error())
		return secretModel, err
	}

	// Convert Json structure to struct model
	json.Unmarshal([]byte(*key.SecretString), &secretModel)

	fmt.Printf("\nSecret name [%s] was read successfully.", secretName)

	return secretModel, nil
}
