// Package discord implements helper function for logging to specific discord server.
package discord

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

// Log log message string to discord via discord webhook.
//
// Logging to discord take some resources and time.
// Recommend that caller only use this function to log error so that it doesn't take to much resources.
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
