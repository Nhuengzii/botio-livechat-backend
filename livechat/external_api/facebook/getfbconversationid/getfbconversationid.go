package getfbconversationid

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Conversations struct {
	Data []*Conversation `json:"data"`
}

type Conversation struct {
	ID string `json:"id"`
}

func GetConversationID(accessToken string, senderID string, pageID string) (_ string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("getfbconversationid.GetConversationID: %w", err)
		}
	}()
	url := fmt.Sprintf("https://graph.facebook.com/v16.0/%v/conversations?platform=Messenger&user_id=%v&access_token=%v",
		pageID, senderID, accessToken)
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
	return conversations.Data[0].ID, nil
}
