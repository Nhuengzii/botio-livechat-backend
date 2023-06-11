package facebook

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ConversationIDs struct {
	Data []Conversation `json:"data"`
}

type Conversation struct {
	ID string `json:"id"`
}

func GetConversationID(accessToken string, senderID string, pageID string) (_ string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("facebook.GetConversationID: %w", err)
		}
	}()
	uri := fmt.Sprintf("https://graph.facebook.com/v16.0/%v/conversations?platform=Messenger&user_id=%v&access_token=%v",
		pageID, senderID, accessToken)

	resp, err := http.Get(uri)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var conversationIDs *ConversationIDs
	err = json.NewDecoder(resp.Body).Decode(conversationIDs)
	if err != nil {
		return "", err
	}

	return conversationIDs.Data[0].ID, nil
}
