package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/storage/amazons3"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func (c *config) handleEvents(ctx context.Context, shopID string, pageID string, bot *linebot.Client, uploader amazons3.Uploader, hookBody *webhookBody) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("handleEvents: %w", err)
		}
	}()
	var stdMessages []*stdmessage.StdMessage
	for _, event := range hookBody.Events {
		switch event.Type {
		case linebot.EventTypeMessage:
			stdMessage, err := c.newStdMessage(shopID, pageID, bot, uploader, event)
			if err != nil {
				if errors.Is(err, errMessageSourceUnsupported) {
					continue
				}
				return err
			}
			stdMessages = append(stdMessages, stdMessage)
		default:
			// TODO implement user join/leave events -> UpdateConversationParticipants
			// info to be updated: name and profile pic
			// ctx is to be passed to db operations here
		}
	}
	for _, stdMessage := range stdMessages {
		stdMessageJSON, err := json.Marshal(stdMessage)
		if err != nil {
			return err
		}
		err = c.snsClient.PublishMessage(c.snsTopicARN, string(stdMessageJSON))
		if err != nil {
			return err
		}
	}
	return nil
}
