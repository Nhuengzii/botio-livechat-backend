package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func logToDiscord(message string) {
	payload := map[string]string{"content": message}
	payloadJSON, _ := json.Marshal(payload)
	_, err := http.Post(discordWebhookURL, "application/json", bytes.NewBuffer(payloadJSON))
	if err != nil {
		log.Println("from logToDiscord: ", err)
	}
}
