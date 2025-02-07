package aws

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"

	"github.com/luisrojas17/ecommerce/models"
)

/*
This function gets the secret configurated in AWS Secret Manager. All data related to secret
are wrapping into SecretModel.

Param "secretName" is a name of the secret to get.
*/
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
