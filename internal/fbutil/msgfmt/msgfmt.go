package msgfmt

import (
	"github.com/Nhuengzii/botio-livechat-backend/internal/fbutil/request"
	"github.com/Nhuengzii/botio-livechat-backend/internal/fbutil/webhook"
	"github.com/Nhuengzii/botio-livechat-backend/pkg/stdmessage"
)

func NewStdMessage(facebookAccessToken string, messageData webhook.MessageData, pageID string) (*stdmessage.StdMessage, error) {
	conversationID, err := request.RequestFacebookConversationID(facebookAccessToken, messageData.Sender.ID, pageID)
	if err != nil {
		return &stdmessage.StdMessage{}, err
	}

	attachments := fmtAttachment(messageData)

	newMessage := stdmessage.StdMessage{
		ShopID:         "1", // TODO:botio API
		Platform:       "Facebook",
		PageID:         pageID,
		ConversationID: conversationID,
		MessageID:      messageData.Message.MessageID,
		Timestamp:      messageData.Timestamp,
		Source: stdmessage.Source{
			UserID:   messageData.Sender.ID,
			UserType: "User",
		},
		Message:     messageData.Message.Text,
		Attachments: attachments,
		ReplyTo: &stdmessage.RepliedMessage{
			MessageID: messageData.Message.ReplyTo.MessageId,
		},
	}

	// check if messageAttachment should be nil
	return &newMessage, nil
}

func fmtAttachment(messageData webhook.MessageData) []*stdmessage.Attachment {
	var attachments []*stdmessage.Attachment
	if len(messageData.Message.Attachments) > 0 {
		for _, attachment := range messageData.Message.Attachments {
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
