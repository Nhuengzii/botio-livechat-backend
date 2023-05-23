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
)

func discordLog(content string) {
	webhookURL := "https://discord.com/api/webhooks/1108750713758175293/U96dYkOWsQYSYrCx6rCFPGvrJ7TY_tMMVmm5IWIAdsCM7ffi_Fa-W9Dfxt7dAd8WNYR2"
	payload := map[string]string{"content": content}
	json_payload, _ := json.Marshal(payload)
	_, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(json_payload))
	if err != nil {
		log.Println("Error sending discord log: ", err)
	}
}
func main() {
	fmt.Println("Hello, World!")
	lambda.Start(Handler)
}

type Attachment struct {
	Type    string `json:"type"`
	Payload struct {
		Src string `json:"src"`
	} `json:"payload"`
}

type Source struct {
	SourceID   string `json:"sourceID"`
	SourceType string `json:"sourceType"`
}

type Message struct {
	ConversationID string       `json:"conversationID"`
	MessageID      string       `json:"messageID"`
	Timestamp      int64        `json:"timeStamp"`
	Source         Source       `json:"source"`
	Message        string       `json:"message"`
	Attachments    []Attachment `json:"attachments"`
}

type BroadcastMessage struct {
	ShopId         string `json:"shopId"`
	Message        string `json:"message"`
	MessageID      string `json:"messageId"`
	ConversationID string `json:"conversationId"`
	Timestamp      int64  `json:"timestamp"`
}

type IncommingMessage struct {
	Action  string  `json:"action"`
	Message Message `json:"message"`
}

func Handler(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	// connectionID := request.RequestContext.ConnectionID
	discordLog("Got broadcast")
	// Unmarshal the message
	var message IncommingMessage
	discordLog("Raw message: " + request.Body)
	err := json.Unmarshal([]byte(request.Body), &message)
	if err != nil {
		discordLog(fmt.Sprint("Error unmarshalling message: ", err))
	}
	discordLog(fmt.Sprintf("Message: %+v", message))

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}
