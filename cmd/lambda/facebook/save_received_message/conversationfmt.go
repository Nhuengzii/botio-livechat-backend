package main

import (
	"context"
	"fmt"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/external_api/facebook/getfbuserprofile"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdconversation"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"
)

func (c *config) newStdConversation(ctx context.Context, message *stdmessage.StdMessage) (_ *stdconversation.StdConversation, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("lambda/facebook/save_received_message/main.newStdConversation: %w", err)
		}
	}()
	facebookCredentials, err := c.dbClient.QueryFacebookPage(ctx, message.PageID)
	if err != nil {
		return nil, err
	}
	userProfile, err := getfbuserprofile.GetUserProfile(facebookCredentials.AccessToken, message.Source.UserID)
	if err != nil {
		return nil, err
	}
	lastActivity, err := message.ToLastActivityString()
	if err != nil {
		return nil, err
	}
	newConversation := &stdconversation.StdConversation{
		ShopID:         message.ShopID,
		Platform:       stdconversation.PlatformFacebook,
		PageID:         message.PageID,
		ConversationID: message.ConversationID,
		ConversationPic: stdconversation.Payload{
			Src: userProfile.ProfilePic,
		},
		UpdatedTime: message.Timestamp,
		Participants: []stdconversation.Participant{
			{
				UserID:   message.Source.UserID,
				Username: userProfile.Name,
				ProfilePic: stdconversation.Payload{
					Src: userProfile.ProfilePic,
				},
			},
		},
		LastActivity: lastActivity,
		IsRead:       false,
	}
	return newConversation, nil
}
