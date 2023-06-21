package main

import (
	"context"
	"encoding/json"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/external_api/instagram/reqigconversationid"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"
)

func (c *config) NewStdMessage(ctx context.Context, messaging Messaging, pageID string) (*stdmessage.StdMessage, error) {
	shop, err := c.dbClient.QueryShop(ctx, pageID)
	if err != nil {
		return nil, err
	}

	instagramCredentials, err := c.dbClient.QueryInstagramAuthentication(ctx, pageID)
	if err != nil {
		return nil, err
	}
	accessToken := instagramCredentials.AccessToken

	var sender stdmessage.UserType
	if messaging.Message.IsEcho {
		sender = stdmessage.UserTypeAdmin
	} else {
		sender = stdmessage.UserTypeUser
	}

	var conversationID string
	if sender == stdmessage.UserTypeUser {
		conversationID, err = reqigconversationid.GetConversationID(accessToken, messaging.Sender.ID, shop.FacebookPageID)
	} else if sender == stdmessage.UserTypeAdmin {
		conversationID, err = reqigconversationid.GetConversationID(accessToken, messaging.Recipient.ID, shop.FacebookPageID)
	}

	if err != nil {
		return nil, err
	}
	attachments, err := fmtAttachment(messaging)
	if err != nil {
		return nil, err
	}
	replyMessage := fmtReplyTo(messaging)

	newMessage := stdmessage.StdMessage{
		ShopID:         shop.ShopID,
		Platform:       stdmessage.PlatformInstagram,
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
		IsDeleted:   messaging.Message.IsDeleted,
	}

	return &newMessage, nil
}

func fmtAttachment(messaging Messaging) ([]stdmessage.Attachment, error) {
	attachments := []stdmessage.Attachment{}
	if len(messaging.Message.Attachments) > 0 {
		for _, attachment := range messaging.Message.Attachments {
			if attachment.AttachmentType != "template" {
				jsonByte, err := json.Marshal(attachment.Payload)
				if err != nil {
					return nil, err
				}
				var basicPayload BasicPayload
				err = json.Unmarshal([]byte(jsonByte), &basicPayload)
				if err != nil {
					return nil, err
				}
				attachments = append(attachments, stdmessage.Attachment{
					AttachmentType: stdmessage.AttachmentType(attachment.AttachmentType),
					Payload:        stdmessage.Payload{Src: basicPayload.Src},
				})
			} else {
				jsonByte, err := json.Marshal(attachment.Payload)
				if err != nil {
					return nil, err
				}
				var templatePayload TemplatePayload
				err = json.Unmarshal(jsonByte, &templatePayload)
				if err != nil {
					return nil, err
				}
				var attachmentType stdmessage.AttachmentType
				var actualJsonPayload []byte
				if len(templatePayload.Generic) != 0 {
					attachmentType = stdmessage.AttachmentTypeIGTemplateGeneric
					actualJsonPayload, err = json.Marshal(templatePayload.Generic)
				} else if len(templatePayload.Product) != 0 {
					attachmentType = stdmessage.AttachmentTypeIGTemplateProduct
					actualJsonPayload, err = json.Marshal(templatePayload.Product)
				} else {
					return nil, errUnknownTemplateType
				}
				if err != nil {
					return nil, err
				}
				attachments = append(attachments, stdmessage.Attachment{
					AttachmentType: attachmentType,
					Payload:        stdmessage.Payload{Src: string(actualJsonPayload)},
				})

			}
		}
	}

	return attachments, nil
}

func fmtReplyTo(messaging Messaging) *stdmessage.RepliedMessage {
	if messaging.Message.ReplyTo.MessageId == "" {
		return nil
	} else {
		return &stdmessage.RepliedMessage{MessageID: messaging.Message.MessageID}
	}
}
