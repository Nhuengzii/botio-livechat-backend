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

	// log detailed info about this event with discordlog
	discordLog(fmt.Sprint("Got event: ", sqsEvent))

	for _, message := range sqsEvent.Records {
		fmt.Println("Got message: ", message.Body)
		body_pretty, _ := json.MarshalIndent(message.Body, "", "  ")
		discordLog(fmt.Sprint("Got message: ", string(body_pretty)))
	}
}
