package conversationfmt

import (
	"fmt"
	"github.com/Nhuengzii/botio-livechat-backend/livechat"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/lineutil/profile"
)

func NewStdConversation(lineChannelAccessToken string, message *livechat.StdMessage) (*livechat.StdConversation, error) {
	userProfile, err := profile.GetUserProfile(lineChannelAccessToken, message.Source.UserID)
	if err != nil {
		return nil, fmt.Errorf("lineutil/conversationfmt.NewStdConversation: %w", err)
	}
	lastActivity, err := message.ToLastActivityString()
	return &livechat.StdConversation{
		ShopID:          message.ShopID,
		PageID:          message.PageID,
		ConversationID:  message.ConversationID,
		ConversationPic: livechat.Payload{Src: userProfile.PictureURL},
		UpdatedTime:     message.Timestamp,
		Participants: []*livechat.Participant{{
			UserID:     message.Source.UserID,
			Username:   userProfile.DisplayName,
			ProfilePic: livechat.Payload{Src: userProfile.PictureURL},
		}},
		LastActivity: lastActivity,
		IsRead:       false,
	}, nil
}
