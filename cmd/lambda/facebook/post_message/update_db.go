package main

import (
	"context"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/external_api/facebook/postfbmessage"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/postmessage"
)

func (c *config) updateDB(ctx context.Context, apiRequestMessage postmessage.Request, fbResponseMessage postfbmessage.SendingMessageResponse, pageID string, conversationID string, psid string) error {
	stdMessage := fmtStdMessage(apiRequestMessage, fbResponseMessage, pageID, conversationID, psid)
	err := c.dbClient.UpdateConversationOnNewMessage(ctx, stdMessage)
	if err != nil {
		return err
	}

	err = c.dbClient.InsertMessage(ctx, stdMessage)
	if err != nil {
		return err
	}
	return nil
}

func fmtStdMessage(apiRequestMessage postmessage.Request, fbResponseMessage postfbmessage.SendingMessageResponse, pageID string, conversationID string, psid string) *stdmessage.StdMessage {
	stdMessage := stdmessage.StdMessage{
		ShopID:         "1",
		Platform:       stdmessage.PlatformFacebook,
		PageID:         pageID,
		ConversationID: conversationID,
		MessageID:      fbResponseMessage.MessageID,
		Timestamp:      fbResponseMessage.Timestamp,
		Source: stdmessage.Source{
			UserID:   pageID, // botio user id?
			UserType: stdmessage.UserTypeAdmin,
		},
		Message: apiRequestMessage.Message,
		Attachments: []*stdmessage.Attachment{
			{
				AttachmentType: stdmessage.AttachmentType(apiRequestMessage.Attachment.AttachmentType),
				Payload:        stdmessage.Payload(apiRequestMessage.Attachment.Payload),
			},
		},
		ReplyTo: &stdmessage.RepliedMessage{
			MessageID: "", // TODO: implement reply
		},
	}

	return &stdMessage
}
