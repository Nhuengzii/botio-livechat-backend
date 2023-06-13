package main

import (
	"encoding/json"
	"fmt"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type webhookBody struct {
	Destination string           `json:"destination"` // bot user id that should receive the webhook
	Events      []*linebot.Event `json:"events"`
}

func parseWebhookBody(body string) (*webhookBody, error) {
	var hookBody webhookBody
	err := json.Unmarshal([]byte(body), &hookBody)
	if err != nil {
		return nil, fmt.Errorf("parseWebhookBody: %w", err)
	}
	return &hookBody, nil
}
