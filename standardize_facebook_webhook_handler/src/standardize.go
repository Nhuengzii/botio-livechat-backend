package main

import "fmt"

// TODO: change value that got hardcoded : ShopID, ConversationID
// TODO: attachments currently support image video audio files only
func Standardize(messageDatas []MessageData, pageID string, standardMessages *[]StandardMessage) {
	for _, messageData := range messageDatas {
		conversationID, err := RequestFacebookConversationID(messageData, pageID)
		if err != nil {
			discordLog(fmt.Sprintf("RequestFacebookConversationID : %v", err))
			return
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
			Attachments: messageData.Message.Attachments,
			ReplyTo: ReplyMessage{
				MessageId: messageData.Message.ReplyTo.MessageId,
			},
		}
		*standardMessages = append(*standardMessages, newMessage)
	}
}
