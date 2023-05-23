package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

type SnsMessage struct {
	Default string `json:"default"`
}

func sendSnsMessage(standardMessage *StandardMessage) error {
	topicArn := os.Getenv("SNS_TOPIC_ARN")
	log.Println(topicArn)
	sess := session.New(&aws.Config{})
	svc := sns.New(sess)
	log.Println("Attempt publishing message from facebook standardizer to SNS")
	log.Printf("%+v", *standardMessage)
	message, err := json.Marshal(*standardMessage)
	if err != nil {
		return err
	}

	snsMessage := SnsMessage{
		Default: string(message),
	}
	snsByte, err := json.Marshal(snsMessage)
	if err != nil {
		return err
	}

	result, err := svc.Publish(&sns.PublishInput{
		MessageStructure: aws.String("json"),
		Message:          aws.String(string(snsByte)),
		TopicArn:         &topicArn,
	})
	if err != nil {
		return err
	}
	log.Println("sns publish result: ", result)
	return nil
}
