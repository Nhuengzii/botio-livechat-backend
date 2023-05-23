package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handle)
}

func handle(ctx context.Context, sqsEvent events.SQSEvent) {
	log.Println("Facebook Message Standardizer handler")
	start := time.Now()
	var recieveMessage RecieveMessage
	var standardMessages []StandardMessage
	for _, record := range sqsEvent.Records {
		err := json.Unmarshal([]byte(record.Body), &recieveMessage)
		if err != nil || recieveMessage.Object != "page" {
			log.Printf("Error unknown webhook object: %v\n", err)
			return
		}
		log.Printf("%+v", recieveMessage)
		for _, message := range recieveMessage.Entry {
			if messaging := message.MessageDatas; len(messaging) != 0 {
				log.Println("messaging field found in recievedMessage")
				// standardize messaging hooks
				err = Standardize(messaging, message.PageID, &standardMessages)
				if err != nil {
					log.Printf("Error standarizing message : %v", err)
					return
				}
			} else {
				log.Printf("Error no message entry : %v", err)
				return
			}
		}
	}

	err := sendSnsMessage(&standardMessages)
	log.Printf("%+v", standardMessages)
	if err != nil {
		log.Println("Error sending SNS message :", err)
		return
	}
	discordLog(fmt.Sprintf("Elapsed: %v", time.Since(start)))
	return
}
