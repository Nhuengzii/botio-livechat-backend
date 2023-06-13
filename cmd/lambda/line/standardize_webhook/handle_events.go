package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func (c *config) handleEvents(ctx context.Context, hookBody *webhookBody) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("handleEvents: %w", err)
		}
	}()
	botUserID := hookBody.Destination
	var stdMessages []*stdmessage.StdMessage
	for _, event := range hookBody.Events {
		switch event.Type {
		case linebot.EventTypeMessage:
			stdMessage, err := c.newStdMessage(ctx, event, botUserID)
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
