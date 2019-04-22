package aeroagg

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
)

func SNSEventRecordToMessage(record *events.SNSEventRecord) (*Message, error) {
	message := record.SNS.Message
	result := &Message{EventKey: &EventKey{}}
	err := json.Unmarshal([]byte(message), result)
	if err == nil {
		result.Date = result.Timestamp.Format("2006-01-02")
	}
	return result, err
}
