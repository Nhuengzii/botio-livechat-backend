package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func discordLog(content string) {
	webhookURL := "https://discord.com/api/webhooks/1108758875068432444/zpT4Tn3bC-q88QdH2XZHTr3POn4vCCajuQdfCks_mXBlKnC3qhKm-KIMclPLGDzdx2vC"
	payload := map[string]string{"content": content}
	json_payload, _ := json.Marshal(payload)
	_, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(json_payload))
	if err != nil {
		log.Println("Error sending discord log: ", err)
	}
}
