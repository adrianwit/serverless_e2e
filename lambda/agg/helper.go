package agg

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)



func SNSEventRecordToMessage(record *events.SNSEventRecord) (*Message, error) {
	message := record.SNS.Message
	result := &Message{EventKey:&EventKey{}}
	err := json.Unmarshal([]byte(message), result)
	if err == nil {
		result.Date = result.Timestamp.Format("2006-01-02")
	}
	return result, err
}


func getDynamoService() (*dynamodb.DynamoDB, error) {
	sess, err :=  session.NewSession()
	if err  != nil {
		return nil, err
	}
	service := dynamodb.New(sess)
	return service, nil
}
