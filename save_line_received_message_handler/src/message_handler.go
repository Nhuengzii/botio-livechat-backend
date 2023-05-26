package main

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

func messageHandler(ctx context.Context, dbc *dbClient, m *botioMessage) error {
	conversationID, err := dbc.checkConversationExists(ctx, m)
	if err != nil {
		// conversation does not exist, create it and insert message
		if err == mongo.ErrNoDocuments {
			c, err := newBotioConversation(m)
			if err != nil {
				return &messageHandlerError{
					message: "couldn't handle botio message",
					err:     err,
				}
			}
			if err := dbc.insertConversation(ctx, c); err != nil {
				return &messageHandlerError{
					message: "couldn't handle botio message",
					err:     err,
				}
			}
			if err := dbc.insertMessage(ctx, m); err != nil {
				return &messageHandlerError{
					message: "couldn't handle botio message",
					err:     err,
				}
			}
			return nil
		} else {
			// some other error
			return &messageHandlerError{
				message: "couldn't handle botio message",
				err:     err,
			}
		}
	}
	// conversation exists, update it and insert message
	if err := dbc.updateConversation(ctx, conversationID, m); err != nil {
		return &messageHandlerError{
			message: "couldn't handle botio message",
			err:     err,
		}
	}
	if err := dbc.insertMessage(ctx, m); err != nil {
		return &messageHandlerError{
			message: "couldn't handle botio message",
			err:     err,
		}
	}
	return nil
}

type messageHandlerError struct {
	message string
	err     error
}

func (e *messageHandlerError) Error() string {
	return e.message + ": " + e.err.Error()
}
