package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func discordLog(content string) {
	webhookURL := "https://discord.com/api/webhooks/1109072106244812870/URR8Y4mCcqqcxF_BlJdl5UTd8KPL7fmwkuJwCXXO_j5w0lkpYT-u__Jxp7fZuBC8oYUs"
	payload := map[string]string{"content": content}
	json_payload, _ := json.Marshal(payload)
	_, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(json_payload))
	if err != nil {
		log.Println("Error sending discord log: ", err)
	}
}
