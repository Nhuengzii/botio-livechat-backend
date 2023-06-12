package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func handleEvents(c *config, hookBody *webhookBody) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("lambda/line/standardize_webhook/main.handleEvents: %w", err)
		}
	}()
	botUserID := hookBody.Destination
	var stdMessages []*stdmessage.StdMessage
	for _, event := range hookBody.Events {
		switch event.Type {
		case linebot.EventTypeMessage:
			stdMessage, err := newStdMessage(event, botUserID)
			if err != nil {
				if errors.Is(err, errMessageSourceUnsupported) {
					continue
				}
				return err
			}
			stdMessages = append(stdMessages, stdMessage)
		default:
			// TODO implement user join/leave events -> updateConversationParticipants
			// info to be updated: group pic, group name, group members, and each member's name and profile pic
			// TODO implement user unsend events -> delete message from db and notify frontend
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
