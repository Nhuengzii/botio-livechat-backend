package main

import (
	"github.com/Nhuengzii/botio-livechat-backend/livechat"
	fbrequest "github.com/Nhuengzii/botio-livechat-backend/livechat/external/facebook"
)

func NewStdMessage(facebookAccessToken string, messaging Messaging, pageID string) (*livechat.StdMessage, error) {
	conversationID, err := fbrequest.RequestFacebookConversationID(facebookAccessToken, messaging.Sender.ID, pageID)
	if err != nil {
		return &livechat.StdMessage{}, err
	}

	attachments := fmtAttachment(messaging)

	newMessage := livechat.StdMessage{
		ShopID:         "1", // TODO:botio API
		Platform:       "Facebook",
		PageID:         pageID,
		ConversationID: conversationID,
		MessageID:      messaging.Message.MessageID,
		Timestamp:      messaging.Timestamp,
		Source: livechat.Source{
			UserID:   messaging.Sender.ID,
			UserType: "User",
		},
		Message:     messaging.Message.Text,
		Attachments: attachments,
		ReplyTo: &livechat.RepliedMessage{
			MessageID: messaging.Message.ReplyTo.MessageId,
		},
	}

	// check if messageAttachment should be nil
	return &newMessage, nil
}

func fmtAttachment(messaging Messaging) []*livechat.Attachment {
	var attachments []*livechat.Attachment
	if len(messaging.Message.Attachments) > 0 {
		for _, attachment := range messaging.Message.Attachments {
			attachments = append(attachments, &livechat.Attachment{
				AttachmentType: livechat.AttachmentType(attachment.AttachmentType),
				Payload:        livechat.Payload{Src: attachment.Payload.Src},
			})
		}
	} else {
		attachments = nil
	}

	return attachments
}
