package main

import (
	"github.com/Nhuengzii/botio-livechat-backend/livechat/external_api/facebook/getfbconversationid"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"
)

func NewStdMessage(messaging Messaging, pageID string) (*stdmessage.StdMessage, error) {
	// TODO: query shopID

	// TODO: query accessToken
	var facebookAccessToken string
	conversationID, err := getfbconversationid.GetConversationID(facebookAccessToken, messaging.Sender.ID, pageID)
	if err != nil {
		return &stdmessage.StdMessage{}, err
	}

	attachments := fmtAttachment(messaging)

	newMessage := stdmessage.StdMessage{
		ShopID:         "1", // TODO:botio API
		Platform:       stdmessage.PlatformFacebook,
		PageID:         pageID,
		ConversationID: conversationID,
		MessageID:      messaging.Message.MessageID,
		Timestamp:      messaging.Timestamp,
		Source: stdmessage.Source{
			UserID:   messaging.Sender.ID,
			UserType: "User",
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
