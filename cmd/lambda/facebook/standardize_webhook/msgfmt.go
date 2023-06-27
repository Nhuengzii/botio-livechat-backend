package main

import (
	"context"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/external_api/facebook/reqfbconversationid"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"
)

func (c *config) NewStdMessage(ctx context.Context, messaging Messaging, pageID string) (*stdmessage.StdMessage, error) {
	shop, err := c.dbClient.QueryShop(ctx, pageID)
	if err != nil {
		return nil, err
	}

	facebookCredentials, err := c.dbClient.QueryFacebookAuthentication(ctx, pageID)
	if err != nil {
		return nil, err
	}
	accessToken := facebookCredentials.AccessToken

	var sender stdmessage.UserType
	if messaging.Message.IsEcho {
		sender = stdmessage.UserTypeAdmin
	} else {
		sender = stdmessage.UserTypeUser
	}

	var conversationID string
	if sender == stdmessage.UserTypeUser {
		conversationID, err = reqfbconversationid.GetConversationID(accessToken, messaging.Sender.ID, pageID)
	} else if sender == stdmessage.UserTypeAdmin {
		conversationID, err = reqfbconversationid.GetConversationID(accessToken, messaging.Recipient.ID, pageID)
	}

	if err != nil {
		return nil, err
	}
	attachments, err := c.fmtAttachment(messaging)
	if err != nil {
		return nil, err
	}
	replyMessage := fmtReplyTo(messaging)

	newMessage := stdmessage.StdMessage{
		ShopID:         shop.ShopID,
		Platform:       stdmessage.PlatformFacebook,
		PageID:         pageID,
		ConversationID: conversationID,
		MessageID:      messaging.Message.MessageID,
		Timestamp:      messaging.Timestamp,
		Source: stdmessage.Source{
			UserID:   messaging.Sender.ID,
			UserType: sender,
		},
		Message:     messaging.Message.Text,
		Attachments: attachments,
		ReplyTo:     replyMessage,
	}

	return &newMessage, nil
}

func fmtReplyTo(messaging Messaging) *stdmessage.RepliedMessage {
	if messaging.Message.ReplyTo.MessageId == "" {
		return nil
	} else {
		return &stdmessage.RepliedMessage{MessageID: messaging.Message.MessageID}
	}
}
