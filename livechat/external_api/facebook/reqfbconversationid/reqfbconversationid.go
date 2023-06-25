// Package reqfbconversationid implement a function to call facebook api request for a conversationID.
//
// # Uses Graph API v16.0
package reqfbconversationid

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// "No conversationID were found after making request"
var errNoConversationIDFound = errors.New("No conversationID were found after making request")

// Conversations is a response body recieved from facebook request which contains a slice of Conversation
type Conversations struct {
	Data []*Conversation `json:"data"` // slice of Conversation
}

// Conversation contain a conversationID string
type Conversation struct {
	ID string `json:"id"` //
}

// GetConversationID makes a faceook API call and returns a string of facebook conversationID,If there is a participants with matching PSID in the conversation.
// Only return conversation in a specific page.
// Return an error if it occurs.
//
// Use facebook page accessToken.
func GetConversationID(accessToken string, psid string, pageID string) (_ string, err error) {
	// important userID is not pageID psid only
	defer func() {
		if err != nil {
			err = fmt.Errorf("reqfbconversationid.GetConversationID: %w", err)
		}
	}()
	url := fmt.Sprintf("https://graph.facebook.com/v16.0/%v/conversations?platform=messenger&user_id=%v&access_token=%v",
		pageID, psid, accessToken)
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
