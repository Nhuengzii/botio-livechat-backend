package main

import (
	"fmt"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/external_api/line/getlineuserprofile"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdconversation"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"
)

func newStdConversation(lineChannelAccessToken string, message *stdmessage.StdMessage) (_ *stdconversation.StdConversation, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("newStdConversation: %w", err)
		}
	}()
	userProfile, err := getlineuserprofile.GetUserProfile(lineChannelAccessToken, message.Source.UserID)
	if err != nil {
		return nil, err
	}
	lastActivity, err := message.ToLastActivityString()
	if err != nil {
		return nil, err
	}
	return &stdconversation.StdConversation{
		ShopID:          message.ShopID,
		Platform:        stdconversation.PlatformLine,
		PageID:          message.PageID,
		ConversationID:  message.ConversationID,
		ConversationPic: stdconversation.Payload{Src: userProfile.PictureURL},
		UpdatedTime:     message.Timestamp,
		Participants: []stdconversation.Participant{
			{
				UserID:     message.Source.UserID,
				Username:   userProfile.DisplayName,
				ProfilePic: stdconversation.Payload{Src: userProfile.PictureURL},
			},
		},
		LastActivity: lastActivity,
		IsRead:       false,
	}, nil
}
