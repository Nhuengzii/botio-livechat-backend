package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/Nhuengzii/botio-livechat-backend/livechat"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/db/mongodb"
)

func updateDB(ctx context.Context, c *config, message *livechat.StdMessage) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("lambda/line/save_received_message/main.updateDB: %w", err)
		}
	}()
	err = c.dbClient.UpdateConversationOnNewMessage(ctx, message)
	if err != nil {
		if errors.Is(err, mongodb.ErrNoConversations) {
			conversation, err := newStdConversation(c.lineChannelAccessToken, message) // TODO get lineChannelAccessToken from db with message.ShopID and message.PageID
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
	err = c.dbClient.InsertMessage(ctx, message)
	if err != nil {
		return err
	}
	return nil
}
