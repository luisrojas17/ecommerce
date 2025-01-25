package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

var Ctx context.Context
var Cfg aws.Config
var err error

/*
This function loads default AWS configuration for this lambda.
Also, The lambda will be executed in North Virginia region.
*/
func Init() {
	Ctx = context.TODO()
	Cfg, err = config.LoadDefaultConfig(Ctx, config.WithDefaultRegion("us-east-1"))

	if err != nil {
		panic("It was an error to load AWS configuration." + err.Error())
	}

}
