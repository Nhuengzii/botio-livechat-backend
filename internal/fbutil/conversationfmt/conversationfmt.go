package conversationfmt

import (
	"github.com/Nhuengzii/botio-livechat-backend/pkg/stdconversation"
	"github.com/Nhuengzii/botio-livechat-backend/pkg/stdmessage"
)

func NewStdConversation(facebookAccessToken string, message *stdmessage.StdMessage) (*stdconversation.StdConversation, error) {
	newConversation := &stdconversation.StdConversation{
		ShopID:          message.ShopID,
		PageID:          message.PageID,
		ConversationID:  message.ConversationID,
		ConversationPic: stdconversation.Payload{
			// Src: ,
		},
		UpdatedTime: message.Timestamp,
		Participants: []*stdconversation.Participant{
			{
				UserID: message.Source.UserID,
				// Username: ,
				ProfilePic: stdconversation.Payload{
					// Src: ,
				},
			},
		},
		LastActivity: message.Message,
		IsRead:       false,
	}
	return newConversation, nil
}
