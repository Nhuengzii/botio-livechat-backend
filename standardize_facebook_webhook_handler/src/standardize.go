package main

// TODO: change value that got hardcoded : ShopID, ConversationID
// TODO: attachments currently support image video audio files only
func Standardize(messageDatas []MessageData, pageID string, standardMessages *[]StandardMessage) {
	for _, messageData := range messageDatas {
		newMessage := StandardMessage{
			ShopID:         "1",
			Platform:       "Facebook",
			PageID:         pageID,
			ConversationID: "t_2422534594589937",
			MessageID:      messageData.Message.MessageID,
			Timestamp:      messageData.Timestamp,
			Source: Source{
				UserID:   messageData.Sender.ID,
				UserType: "user",
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
