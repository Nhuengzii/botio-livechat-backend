package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func discordLog(content string) {
	webhookURL := "https://discord.com/api/webhooks/1109037432814452787/yzkc5Qohr7cK75k2Cb-SjoMml7mZpwXiSIQxN8qG_MFwsM4Nr30Sp6S4ofL_f3Pxhm0d"
	payload := map[string]string{"content": content}
	json_payload, _ := json.Marshal(payload)
	_, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(json_payload))
	if err != nil {
		log.Println("Error sending discord log: ", err)
	}
}
