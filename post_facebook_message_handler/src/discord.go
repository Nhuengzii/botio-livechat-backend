package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func discordLog(content string) {
	webhookURL := "https://discord.com/api/webhooks/1110130057613148221/Q6W07r2X6JU5b9h6VBZhqsrYIac3OsHOGBbCLRiyjabVsN2yOh7lbzmnKag16A1b7ovf"
	payload := map[string]string{"content": content}
	json_payload, _ := json.Marshal(payload)
	_, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(json_payload))
	if err != nil {
		log.Println("Error sending discord log: ", err)
	}
}
