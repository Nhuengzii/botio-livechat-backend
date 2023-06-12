package main

import (
	"context"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/external_api/facebook/getfbconversationid"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"
)

func (c *config) NewStdMessage(ctx context.Context, messaging Messaging, pageID string) (*stdmessage.StdMessage, error) {
	shop, err := c.DbClient.QueryShop(ctx, pageID)
	if err != nil {
		return nil, err
	}

	facebookCredentials, err := c.DbClient.QueryFacebookPageCredentials(ctx, pageID)
	if err != nil {
		return nil, err
	}
	accessToken := facebookCredentials.AccessToken
	conversationID, err := getfbconversationid.GetConversationID(accessToken, messaging.Sender.ID, pageID)
	if err != nil {
		return nil, err
	}

	attachments := fmtAttachment(messaging)

	newMessage := stdmessage.StdMessage{
		ShopID:         shop.ShopID,
		Platform:       stdmessage.PlatformFacebook,
		PageID:         pageID,
		ConversationID: conversationID,
		MessageID:      messaging.Message.MessageID,
		Timestamp:      messaging.Timestamp,
		Source: stdmessage.Source{
			UserID:   messaging.Sender.ID,
			UserType: stdmessage.UserTypeUser,
		},
		Message:     messaging.Message.Text,
		Attachments: attachments,
		ReplyTo: &stdmessage.RepliedMessage{
			MessageID: messaging.Message.ReplyTo.MessageId,
		},
	}

	return &newMessage, nil
}

func fmtAttachment(messaging Messaging) []*stdmessage.Attachment {
	var attachments []*stdmessage.Attachment
	if len(messaging.Message.Attachments) > 0 {
		for _, attachment := range messaging.Message.Attachments {
			attachments = append(attachments, &stdmessage.Attachment{
				AttachmentType: stdmessage.AttachmentType(attachment.AttachmentType),
				Payload:        stdmessage.Payload{Src: attachment.Payload.Src},
			})
		}
	} else {
		attachments = nil
	}

	return attachments
}
