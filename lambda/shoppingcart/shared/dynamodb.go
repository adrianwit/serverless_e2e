package shared

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var dynamoDb *dynamodb.DynamoDB

func GetDynamoService() (*dynamodb.DynamoDB, error) {
	if dynamoDb != nil {
		return dynamoDb, nil
	}
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	dynamoDb = dynamodb.New(sess)
	return dynamoDb, nil
}

