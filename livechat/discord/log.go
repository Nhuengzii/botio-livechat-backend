// Package discord includes a helper function for logging to a specific discord server.
package discord

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

// Log logs a message string to discord via discord webhook.
// Logging to discord takes some time, so it is recommended to only use this function to log errors.
// If webhookURL is empty, the message will not be logged.
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
