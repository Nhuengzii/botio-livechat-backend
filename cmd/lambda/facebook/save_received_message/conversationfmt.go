package main

import (
	"context"
	"fmt"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/external_api/facebook/reqfbuserprofile"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/external_api/facebook/reqfbuserpsid"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdconversation"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"
)

func (c *config) newStdConversation(ctx context.Context, message *stdmessage.StdMessage) (_ *stdconversation.StdConversation, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("lambda/facebook/save_received_message/main.newStdConversation: %w", err)
		}
	}()
	facebookCredentials, err := c.dbClient.QueryFacebookAuthentication(ctx, message.PageID)
	if err != nil {
		return nil, err
	}

	var userProfile *reqfbuserprofile.UserProfile
	var userID string
	if message.Source.UserType == stdmessage.UserTypeUser {
		// no need to query user's psid
		userProfile, err = reqfbuserprofile.GetUserProfile(facebookCredentials.AccessToken, message.Source.UserID)
		if err != nil {
			return nil, err
		}
		userID = message.Source.UserID
	} else if message.Source.UserType == stdmessage.UserTypeAdmin {
		// query for user's psid from pageID
		psid, err := reqfbuserpsid.GetUserPSID(facebookCredentials.AccessToken, message.PageID, message.ConversationID)
		if err != nil {
			return nil, err
		} else if psid == "" {
			return nil, errCannotGetUserPSID
		}

		userProfile, err = reqfbuserprofile.GetUserProfile(facebookCredentials.AccessToken, psid)
		if err != nil {
			return nil, err
		}

		userID = psid
	} else {
		return nil, errUnsupportedUserType
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
				UserID:   userID,
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
