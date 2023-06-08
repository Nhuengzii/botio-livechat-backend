package request

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func RequestFacebookConversationID(access_token string, senderID string, pageID string) (string, error) {
	uri := fmt.Sprintf("https://graph.facebook.com/v16.0/%v/conversations?platform=Messenger&user_id=%v&access_token=%v",
		pageID, senderID, access_token)

	resp, err := http.Get(uri)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var body ResponseConversationID
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return "", err
	}

	return body.Data[0].Id, nil
}

type ResponseConversationID struct {
	Data []Conversation `json:"data"`
}

type Conversation struct {
	Id string `json:"id"`
}
