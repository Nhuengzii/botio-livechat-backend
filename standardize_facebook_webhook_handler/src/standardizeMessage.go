package main

import "fmt"

// TODO: change value that got hardcoded : ShopID, ConversationID
// TODO: attachments currently support image video audio files only
func Standardize(messageData MessageData, pageID string, standardMessages *[]StandardMessage) error {
	conversationID, err := RequestFacebookConversationID(messageData, pageID)
	if err != nil {
		discordLog(fmt.Sprintf("RequestFacebookConversationID : %v", err))
		return err
	}
	newMessage := StandardMessage{
		ShopID:         "1", // TODO:botio API
		Platform:       "Facebook",
		PageID:         pageID,
		ConversationID: conversationID, // TODO: get from fb?
		MessageID:      messageData.Message.MessageID,
		Timestamp:      messageData.Timestamp,
		Source: Source{
			UserID:   messageData.Sender.ID,
			UserType: "User",
		},
		Message:     messageData.Message.Text,
		Attachments: []AttachmentOutput{},
		ReplyTo: ReplyMessage{
			MessageId: messageData.Message.ReplyTo.MessageId,
		},
	}
	for _, attachment := range messageData.Message.Attachments {
		newMessage.Attachments = append(newMessage.Attachments, AttachmentOutput{
			AttachmentType: attachment.AttachmentType,
			Payload:        PayloadTypeOutput{Src: attachment.Payload.Src},
		})
	}
	*standardMessages = append(*standardMessages, newMessage)

	return nil
}
