package msgfmt

import (
	"github.com/Nhuengzii/botio-livechat-backend/livechat"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/fbutil/request"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/fbutil/webhook"
)

func NewStdMessage(facebookAccessToken string, messageData webhook.MessageData, pageID string) (*livechat.StdMessage, error) {
	conversationID, err := request.RequestFacebookConversationID(facebookAccessToken, messageData.Sender.ID, pageID)
	if err != nil {
		return &livechat.StdMessage{}, err
	}

	attachments := fmtAttachment(messageData)

	newMessage := livechat.StdMessage{
		ShopID:         "1", // TODO:botio API
		Platform:       "Facebook",
		PageID:         pageID,
		ConversationID: conversationID,
		MessageID:      messageData.Message.MessageID,
		Timestamp:      messageData.Timestamp,
		Source: livechat.Source{
			UserID:   messageData.Sender.ID,
			UserType: "User",
		},
		Message:     messageData.Message.Text,
		Attachments: attachments,
		ReplyTo: &livechat.RepliedMessage{
			MessageID: messageData.Message.ReplyTo.MessageId,
		},
	}

	// check if messageAttachment should be nil
	return &newMessage, nil
}

func fmtAttachment(messageData webhook.MessageData) []*livechat.Attachment {
	var attachments []*livechat.Attachment
	if len(messageData.Message.Attachments) > 0 {
		for _, attachment := range messageData.Message.Attachments {
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
