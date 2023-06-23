package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/line/line-bot-sdk-go/v7/linebot"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/db/mongodb"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"
)

func (c *config) handleMessage(ctx context.Context, bot *linebot.Client, message stdmessage.StdMessage) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("handleMessage: %w", err)
		}
	}()
	err = c.dbClient.UpdateConversationOnNewMessage(ctx, &message)
	if err != nil {
		if errors.Is(err, mongodb.ErrNoDocuments) {
			conversation, err := newStdConversation(bot, &message)
			if err != nil {
				return err
			}
			err = c.dbClient.InsertConversation(ctx, conversation)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	err = c.dbClient.InsertMessage(ctx, &message)
	if err != nil {
		return err
	}
	return nil
}
