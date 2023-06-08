package discord

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func Log(webhookURL string, message string) {
	if webhookURL != "" {
		payload := map[string]string{"content": message}
		payloadJSON, _ := json.Marshal(payload)
		_, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payloadJSON))
		if err != nil {
			log.Println("discord.Log: ", err)
		}
	}
}
