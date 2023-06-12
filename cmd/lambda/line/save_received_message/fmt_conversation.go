package main

import (
	"fmt"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdconversation"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"
)

func newStdConversation(lineChannelAccessToken string, message *stdmessage.StdMessage) (*stdconversation.StdConversation, error) {
	userProfile, err := getUserProfile(lineChannelAccessToken, message.Source.UserID)
	if err != nil {
		return nil, fmt.Errorf("lineutil/conversationfmt.NewStdConversation: %w", err)
	}
	lastActivity, err := message.ToLastActivityString()
	return &stdconversation.StdConversation{
		ShopID:          message.ShopID,
		PageID:          message.PageID,
		Platform:        stdconversation.PlatformLine,
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
