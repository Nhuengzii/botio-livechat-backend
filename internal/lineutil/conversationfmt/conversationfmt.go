package conversationfmt

import (
	"fmt"
	"github.com/Nhuengzii/botio-livechat-backend/internal/lineutil/profile"
	"github.com/Nhuengzii/botio-livechat-backend/pkg/stdconversation"
	"github.com/Nhuengzii/botio-livechat-backend/pkg/stdmessage"
)

func NewStdConversation(lineChannelAccessToken string, message *stdmessage.StdMessage) (*stdconversation.StdConversation, error) {
	userProfile, err := profile.GetUserProfile(lineChannelAccessToken, message.Source.UserID)
	if err != nil {
		return nil, fmt.Errorf("lineutil/conversationfmt.NewStdConversation: %w", err)
	}
	lastActivity, err := message.ToLastActivityString()
	return &stdconversation.StdConversation{
		ShopID:          message.ShopID,
		PageID:          message.PageID,
		ConversationID:  message.ConversationID,
		ConversationPic: stdconversation.Payload{Src: userProfile.PictureURL},
		UpdatedTime:     message.Timestamp,
		Participants: []*stdconversation.Participant{{
			UserID:     message.Source.UserID,
			Username:   userProfile.DisplayName,
			ProfilePic: stdconversation.Payload{Src: userProfile.PictureURL},
		}},
		LastActivity: lastActivity,
		IsRead:       false,
	}, nil
}
