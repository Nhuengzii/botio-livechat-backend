package conversationfmt

import (
	"github.com/Nhuengzii/botio-livechat-backend/livechat"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/fbutil/request"
)

func NewStdConversation(facebookAccessToken string, message *livechat.StdMessage) (*livechat.StdConversation, error) {
	userProfile, err := request.RequestFacebookUserProfile(facebookAccessToken, message.Source.UserID)
	if err != nil {
		return &livechat.StdConversation{}, err
	}
	newConversation := &livechat.StdConversation{
		ShopID:         message.ShopID,
		PageID:         message.PageID,
		ConversationID: message.ConversationID,
		ConversationPic: livechat.Payload{
			Src: userProfile.ProfilePic,
		},
		UpdatedTime: message.Timestamp,
		Participants: []*livechat.Participant{
			{
				UserID:   message.Source.UserID,
				Username: userProfile.Name,
				ProfilePic: livechat.Payload{
					Src: userProfile.ProfilePic,
				},
			},
		},
		LastActivity: message.Message,
		IsRead:       false,
	}
	return newConversation, nil
}
