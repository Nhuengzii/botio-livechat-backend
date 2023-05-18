package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-redis/redis/v8"
)

type Source struct {
	UserID   string `json:"userID"`
	UserType string `json:"userType"`
}
type PayloadType struct {
	Src string `json:"url"`
}
type Attachment struct {
	AttachmentType string      `json:"type"`
	Payload        PayloadType `json:"payload"`
}
type ReplyMessage struct {
	MessageId string `json:"messageID"`
}
type StandardMessage struct {
	ShopID         string       `json:"shopID"`
	PageID         string       `json:"pageID"`
	ConversationID string       `json:"conversationID"`
	MessageID      string       `json:"messageID"`
	Timestamp      int64        `json:"timestamp"`
	Source         Source       `json:"source"`
	Message        string       `json:"message"`
	Attachments    []Attachment `json:"attachments"`
	ReplyTo        ReplyMessage `json:"replyTo"`
}

func main() {
	fmt.Println("Hello Me")
	lambda.Start(Handler)
}
func discordLog(content string) {
	webhookURL := "https://discord.com/api/webhooks/1108750713758175293/U96dYkOWsQYSYrCx6rCFPGvrJ7TY_tMMVmm5IWIAdsCM7ffi_Fa-W9Dfxt7dAd8WNYR2"
	payload := map[string]string{"content": content}
	json_payload, _ := json.Marshal(payload)
	_, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(json_payload))
	if err != nil {
		log.Println("Error sending discord log: ", err)
	}
}

func Handler(ctx context.Context, sqsEvent events.SQSEvent) {
	fmt.Println("Hello Me")
	discordLog(fmt.Sprint("Got connect: ", sqsEvent.Records))
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis-15520.c252.ap-southeast-1-1.ec2.cloud.redislabs.com:15520",
		Password: "dcesPhFIPwWrItb2yaNe5UT0sbhv9FJk",
	})
	fmt.Println(rdb)

	for _, message := range sqsEvent.Records {
		fmt.Printf("The message %s for event source %s = %s \n", message.MessageId, message.EventSource, message.Body)
		var standardMessage StandardMessage
		json.Unmarshal([]byte(message.Body), &standardMessage)
		discordLog(fmt.Sprint("Got message: ", standardMessage.Message))
	}

}
