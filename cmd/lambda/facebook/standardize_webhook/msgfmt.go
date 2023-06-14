package main

import (
	"context"
	"encoding/json"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/external_api/facebook/getfbconversationid"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"
)

func (c *config) NewStdMessage(ctx context.Context, messaging Messaging, pageID string) (*stdmessage.StdMessage, error) {
	shop, err := c.dbClient.QueryShop(ctx, pageID)
	if err != nil {
		return nil, err
	}

	facebookCredentials, err := c.dbClient.QueryFacebookPage(ctx, pageID)
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
		conversationID, err = getfbconversationid.GetConversationID(accessToken, messaging.Sender.ID, pageID)
	} else if sender == stdmessage.UserTypeAdmin {
		conversationID, err = getfbconversationid.GetConversationID(accessToken, messaging.Recipient.ID, pageID)
	}

	if err != nil {
		return nil, err
	}
	attachments, err := fmtAttachment(messaging)
	if err != nil {
		return nil, err
	}

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
		ReplyTo: &stdmessage.RepliedMessage{
			MessageID: messaging.Message.ReplyTo.MessageId,
		},
	}

	return &newMessage, nil
}

func fmtAttachment(messaging Messaging) ([]*stdmessage.Attachment, error) {
	var attachments []*stdmessage.Attachment
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
				attachments = append(attachments, &stdmessage.Attachment{
					AttachmentType: stdmessage.AttachmentType(attachment.AttachmentType),
					Payload:        stdmessage.Payload{Src: basicPayload.Src},
				})
			} else {
				jsonByte, err := json.Marshal(attachment.Payload) // actual payload
				if err != nil {
					return nil, err
				}
				var templatePayload TemplatePayload
				err = json.Unmarshal(jsonByte, &templatePayload)
				if err != nil {
					return nil, err
				}
				var attachmentType stdmessage.AttachmentType
				if templatePayload.TemplateType == "button" {
					attachmentType = stdmessage.AttachmentTypeFBTemplateButton
				} else if templatePayload.TemplateType == "coupon" {
					attachmentType = stdmessage.AttachmentTypeFBTemplateCoupon
				} else if templatePayload.TemplateType == "customer_feedback" {
					attachmentType = stdmessage.AttachmentTypeFBTemplateCustomerFeedback
				} else if templatePayload.TemplateType == "generic" {
					attachmentType = stdmessage.AttachmentTypeFBTemplateGeneric
				} else if templatePayload.TemplateType == "media" {
					attachmentType = stdmessage.AttachmentTypeFBTemplateMedia
				} else if templatePayload.TemplateType == "product" {
					attachmentType = stdmessage.AttachmentTypeFBTemplateProduct
				} else if templatePayload.TemplateType == "receipt" {
					attachmentType = stdmessage.AttachmentTypeFBTemplateReceipt
				} else if templatePayload.TemplateType == "customer_information" {
					attachmentType = stdmessage.AttachmentTypeFBTemplateStructuredInformation
				} else {
					return nil, errUnknownTemplateType
				}
				attachments = append(attachments, &stdmessage.Attachment{
					AttachmentType: attachmentType,
					Payload:        stdmessage.Payload{Src: string(jsonByte)},
				})

			}
		}
	} else {
		attachments = nil
	}

	return attachments, nil
}
