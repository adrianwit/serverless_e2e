package agg

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Service interface {
	Consume(message *Message ) error
}

type service struct {}

func (s *service) Consume(message *Message ) error  {
	dyndbService, err := getDynamoService()
	if err != nil {
		return err
	}
	table := aws.String("Events")
	key, err := dynamodbattribute.MarshalMap(message.EventKey)
	if err != nil {
		return fmt.Errorf("unable to marshal key %v", err)
	}
	if getOutput, err := dyndbService.GetItem(&dynamodb.GetItemInput{
		TableName:table,
		Key:key,
	});err != nil || len(getOutput.Item) == 0  {
		event := &Event{
			EventType:message.EventType,
			Date:message.Date,
			Quantity:0,
		}
		item, err := dynamodbattribute.MarshalMap(event)
		if err != nil {
			return err
		}
		if _, err = dyndbService.PutItem(&dynamodb.PutItemInput{
			TableName:table,
			Item:item,
		});err != nil {
			return err
		}
	}
	delta, err := dynamodbattribute.MarshalMap(Counter{Value: 1,})
	if err != nil {
		return fmt.Errorf("unable to marshal counter %v", err)
	}
	input := &dynamodb.UpdateItemInput{
		TableName:                 table,
		Key:                       key,
		UpdateExpression:          aws.String("set Qty = Qty + :Delta"),
		ExpressionAttributeValues: delta,
		ReturnValues:              aws.String("NONE"),
	}
	_, err = dyndbService.UpdateItem(input)
	return err
}


//New returns new service
func New() Service {
	return &service{}
}



type Counter struct {
	Value int `json:":Delta"`
}