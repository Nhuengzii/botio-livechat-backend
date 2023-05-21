package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

var discordWebhookUrl = os.Getenv("DISCORD_WEBHOOK_URL")

func discordLog(content string) {
	webhookURL := discordWebhookUrl
	payload := map[string]string{"content": content}
	json_payload, _ := json.Marshal(payload)
	_, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(json_payload))
	if err != nil {
		log.Println("Error sending discord log: ", err)
	}
}
