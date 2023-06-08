package webhook

import (
	"github.com/Nhuengzii/botio-livechat-backend/internal/fbutil/request"
	"github.com/Nhuengzii/botio-livechat-backend/pkg/stdmessage"
)

func (messageData MessageData) StandardizeMessage(facebookAccessToken string, pageID string, standardMessage *stdmessage.StdMessage) error {
	conversationID, err := request.RequestFacebookConversationID(facebookAccessToken, messageData.Sender.ID, pageID)
	if err != nil {
		return err
	}
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
		Attachments: []*stdmessage.Attachment{},
		ReplyTo: &stdmessage.RepliedMessage{
			MessageID: messageData.Message.ReplyTo.MessageId,
		},
	}
	for _, attachment := range messageData.Message.Attachments {
		newMessage.Attachments = append(newMessage.Attachments, &stdmessage.Attachment{
			AttachmentType: stdmessage.AttachmentType(attachment.AttachmentType),
			Payload:        stdmessage.Payload{Src: attachment.Payload.Src},
		})
	}
	*standardMessage = newMessage

	return nil
}
