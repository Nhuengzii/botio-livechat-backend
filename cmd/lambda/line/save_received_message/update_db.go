package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/db/mongodb"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"
)

func (c *config) updateDB(ctx context.Context, message *stdmessage.StdMessage) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("lambda/line/save_received_message/main.updateDB: %w", err)
		}
	}()
	err = c.dbClient.UpdateConversationOnNewMessage(ctx, message)
	if err != nil {
		if errors.Is(err, mongodb.ErrNoDocuments) {
			conversation, err := newStdConversation(c.lineChannelAccessToken, message) // TODO get lineChannelAccessToken from caller
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
