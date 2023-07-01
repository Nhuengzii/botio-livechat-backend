// Package reqigconversationid implement a function to call instagram api request for a conversationID.
//
// # Uses Graph API v16.0
package reqigconversationid

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// "No conversationID were found after making request"
var errNoConversationIDFound = errors.New("No conversationID found ")

// Conversations is a response body recieved from graph API request which contains a slice of Conversation
type Conversations struct {
	Data []*Conversation `json:"data"`
}

// Conversation contain a conversationID string
type Conversation struct {
	ID string `json:"id"`
}

// GetConversationID makes a graph API call and returns a string of instagram conversationID,If there is a participants with matching PSID in the conversation.
// Only return conversation in a specify page.
// Return an error if it occurs.
//
// Use instagram page accessToken.
func GetConversationID(accessToken string, igsid string, facebookPageID string) (_ string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("reqigconversationid.GetConversationID: %w", err)
		}
	}()
	url := fmt.Sprintf("https://graph.facebook.com/v16.0/%v/conversations?platform=instagram&user_id=%v&access_token=%v",
		facebookPageID, igsid, accessToken)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var conversations Conversations
	err = json.NewDecoder(resp.Body).Decode(&conversations)
	if err != nil {
		return "", err
	}
	if len(conversations.Data) <= 0 {
		return "", errNoConversationIDFound
	}
	return conversations.Data[0].ID, nil
}
