package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/discord"
)

func (c *config) handleRecieveWebhook(ctx context.Context, recieveWebhook *ReceiveWebhook) error {
	if recieveWebhook.Object != "page" {
		return errUnknownWebhookObject
	}
	for _, entry := range recieveWebhook.Entries {
		err := c.handleWebhookEntry(ctx, &entry)
		if err != nil {
			discord.Log(c.DiscordWebhookURL, fmt.Sprintf("error handling webhook entry : %v", err))
			return err
		}
	}
	return nil
}

func (c *config) handleWebhookEntry(ctx context.Context, message *Entry) error {
	if len(message.Messagings) <= 0 {
		return errNoMessageEntry
	}

	for _, messaging := range message.Messagings {
		if messaging.Message.MessageID != "" {
			// standardize messaging hooks
			var standardMessage *stdmessage.StdMessage
			standardMessage, err := c.NewStdMessage(ctx, messaging, message.PageID)
			if err != nil {
				return err
			}

			standardMessageJSON, err := json.Marshal(standardMessage)
			if err != nil {
				return err
			}
			err = c.SnsClient.PublishMessage(c.SnsTopicARN, string(standardMessageJSON))
			if err != nil {
				return err
			}
		} else {
			return errUnknownWebhookType
		}
	}
	return nil
}
