package main

import (
	"context"

	"github.com/Nhuengzii/botio-livechat-backend/livechat"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/request/postmessagreq"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/external/fbrequest"
)

func (c *config) updateDB(ctx context.Context, apiRequestMessage postmessagreq.Request, fbResponseMessage fbrequest.FBSendMsgResponse, pageID string, conversationID string, psid string) error {
	stdMessage := fmtStdMessage(apiRequestMessage, fbResponseMessage, pageID, conversationID, psid)
	err := c.DbClient.UpdateConversationOnNewMessage(ctx, stdMessage)
	if err != nil {
		return err
	}

	err = c.DbClient.InsertMessage(ctx, stdMessage)
	if err != nil {
		return err
	}
	return nil
}

func fmtStdMessage(apiRequestMessage postmessagreq.Request, fbResponseMessage fbrequest.FBSendMsgResponse, pageID string, conversationID string, psid string) *livechat.StdMessage {
	stdMessage := livechat.StdMessage{
		ShopID:         "1",
		Platform:       "Facebook",
		PageID:         pageID,
		ConversationID: conversationID,
		MessageID:      fbResponseMessage.MessageID,
		Timestamp:      fbResponseMessage.Timestamp,
		Source: livechat.Source{
			UserID:   pageID, // botio user id?
			UserType: "Admin",
		},
		Message: apiRequestMessage.Message,
		Attachments: []*livechat.Attachment{
			{
				AttachmentType: livechat.AttachmentType(apiRequestMessage.Attachment.AttachmentType),
				Payload:        livechat.Payload(apiRequestMessage.Attachment.Payload),
			},
		},
		ReplyTo: &livechat.RepliedMessage{
			MessageID: "", // TODO: implement reply
		},
	}

	return &stdMessage
}
