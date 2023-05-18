package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handle)
}

func handle(ctx context.Context, sqsEvent events.SQSEvent) {
	log.Println("Facebook Message Standardizer handler")
	var recieveMessage RecieveMessage
	var standardMessages []StandardMessage
	for _, record := range sqsEvent.Records {
		err := json.Unmarshal([]byte(record.Body), &recieveMessage)
		if err != nil {
			log.Printf("Error unmarshal Record.Body: %v\n", err)
			return
		}
		log.Printf("%+v", recieveMessage)
		for _, message := range recieveMessage.Entry {
			if messaging := message.MessageDatas; messaging != nil {
				log.Println("messaging field found in recievedMessage")
				// standardize messaging hooks
				Standardize(messaging, message.PageID, &standardMessages)
			}
		}
	}

	err := sendSnsMessage(&standardMessages)
	log.Printf("%+v", standardMessages)
	if err != nil {
		log.Println("Error sending SNS message :", err)
	}
	return
}
