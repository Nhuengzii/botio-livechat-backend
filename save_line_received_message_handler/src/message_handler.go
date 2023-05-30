package main

import (
	"context"
	"fmt"
)

func messageHandler(ctx context.Context, dbc *dbClient, m *botioMessage) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("messageHandler: %w", err)
		}
	}()
	exists, err := dbc.checkConversationExists(ctx, m)
	// some unexpected error
	if err != nil {
		return err
	}
	// no conversation exists; create a new one and insert the message
	if !exists {
		conversation, err := newBotioConversation(m)
		if err != nil {
			return err
		}
		if err := dbc.insertConversation(ctx, conversation); err != nil {
			return err
		}
		if err := dbc.insertMessage(ctx, m); err != nil {
			return err
		}
	} else {
		// conversation exists; update the conversation and insert the message
		if err := dbc.updateConversationOfMessage(ctx, m); err != nil {
			return err
		}
		if err := dbc.insertMessage(ctx, m); err != nil {
			return err
		}
	}
	return nil
}
