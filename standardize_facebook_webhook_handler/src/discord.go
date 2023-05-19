package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func discordLog(content string) {
	webhookURL := "https://discord.com/api/webhooks/1109008647394181200/q-qCLP6-LwmzSdfJDKqLqaA0_mBFo5b8BQ1qk1WE3oeJweQ1FoiQKjSA_tQchpvrHo1S"
	payload := map[string]string{"content": content}
	json_payload, _ := json.Marshal(payload)
	_, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(json_payload))
	if err != nil {
		log.Println("Error sending discord log: ", err)
	}
}
