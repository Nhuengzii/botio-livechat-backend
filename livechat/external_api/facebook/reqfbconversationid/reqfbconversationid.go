package reqfbconversationid

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var errNoConversationIDFound = errors.New("No conversationID found ")

type Conversations struct {
	Data []*Conversation `json:"data"`
}

type Conversation struct {
	ID string `json:"id"`
}

func GetConversationID(accessToken string, psid string, pageID string) (_ string, err error) {
	// important userID is not pageID psid only
	defer func() {
		if err != nil {
			err = fmt.Errorf("reqfbconversationid.GetConversationID: %w", err)
		}
	}()
	url := fmt.Sprintf("https://graph.facebook.com/v16.0/%v/conversations?platform=Messenger&user_id=%v&access_token=%v",
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