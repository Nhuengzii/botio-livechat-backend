package main

import (
	"context"
	"log"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/external_api/facebook/postfbmessage"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/postmessage"
)

func (c *config) updateDB(ctx context.Context, apiRequestMessage postmessage.Request, fbResponseMessage postfbmessage.SendingMessageResponse, shopID string, pageID string, conversationID string, psid string) error {
	stdMessage := fmtStdMessage(apiRequestMessage, fbResponseMessage, shopID, pageID, conversationID, psid)
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

func fmtStdMessage(apiRequestMessage postmessage.Request, fbResponseMessage postfbmessage.SendingMessageResponse, shopID string, pageID string, conversationID string, psid string) *stdmessage.StdMessage {
	log.Println(apiRequestMessage.Attachment.AttachmentType)
	log.Println(apiRequestMessage.Attachment.Payload)
	var attachments []*stdmessage.Attachment
	if apiRequestMessage.Attachment.AttachmentType != "" {
		attachments = append(attachments, &stdmessage.Attachment{
			AttachmentType: stdmessage.AttachmentType(apiRequestMessage.Attachment.AttachmentType),
			Payload:        stdmessage.Payload(apiRequestMessage.Attachment.Payload),
		})
	}
	stdMessage := stdmessage.StdMessage{
		ShopID:         shopID,
		Platform:       stdmessage.PlatformFacebook,
		PageID:         pageID,
		ConversationID: conversationID,
		MessageID:      fbResponseMessage.MessageID,
		Timestamp:      fbResponseMessage.Timestamp,
		Source: stdmessage.Source{
			UserID:   pageID, // botio user id?
			UserType: stdmessage.UserTypeAdmin,
		},
		Message:     apiRequestMessage.Message,
		Attachments: attachments,
		ReplyTo: &stdmessage.RepliedMessage{
			MessageID: "", // TODO: implement reply
		},
	}

	return &stdMessage
}
