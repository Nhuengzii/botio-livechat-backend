package conversationfmt

import (
	"github.com/Nhuengzii/botio-livechat-backend/internal/fbutil/request"
	"github.com/Nhuengzii/botio-livechat-backend/pkg/stdconversation"
	"github.com/Nhuengzii/botio-livechat-backend/pkg/stdmessage"
)

func NewStdConversation(facebookAccessToken string, message *stdmessage.StdMessage) (*stdconversation.StdConversation, error) {
	userProfile, err := request.RequestFacebookUserProfile(facebookAccessToken, message.Source.UserID)
	if err != nil {
		return &stdconversation.StdConversation{}, err
	}
	newConversation := &stdconversation.StdConversation{
		ShopID:         message.ShopID,
		PageID:         message.PageID,
		ConversationID: message.ConversationID,
		ConversationPic: stdconversation.Payload{
			Src: userProfile.ProfilePic,
		},
		UpdatedTime: message.Timestamp,
		Participants: []*stdconversation.Participant{
			{
				UserID:   message.Source.UserID,
				Username: userProfile.Name,
				ProfilePic: stdconversation.Payload{
					Src: userProfile.ProfilePic,
				},
			},
		},
		LastActivity: message.Message,
		IsRead:       false,
	}
	return newConversation, nil
}
