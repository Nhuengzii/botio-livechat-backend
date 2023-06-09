package main

import (
	"fmt"

	"github.com/Nhuengzii/botio-livechat-backend/livechat"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/discord"
)

func (c *config) handleRecieveWebhook(recieveWebhook *ReceiveWebhook) error {
	if recieveWebhook.Object != "page" {
		return errUnknownWebhookObject
	}
	for _, entry := range recieveWebhook.Entries {
		err := c.handleWebhookEntry(&entry)
		if err != nil {
			discord.Log(c.DiscordWebhookURL, fmt.Sprintf("error handling webhook entry : %v", err))
			return err
		}
	}
	return nil
}

func (c *config) handleWebhookEntry(message *Entry) error {
	if len(message.Messagings) <= 0 {
		return errNoMessageEntry
	}

	for _, messageData := range message.Messagings {
		if messageData.Message.MessageID != "" {
			// standardize messaging hooks
			var standardMessage *livechat.StdMessage
			standardMessage, err := msgfmt.NewStdMessage(c.FacebookPageAccessToken, messageData, message.PageID)
			if err != nil {
				return err
			}

			err = c.SnsClient.PublishMessage(c.SnsQueueURL, *standardMessage)
			if err != nil {
				return err
			}
		} else {
			return errUnknownWebhookType
		}
	}
	return nil
}
