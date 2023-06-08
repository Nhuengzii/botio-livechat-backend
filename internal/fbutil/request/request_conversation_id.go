package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func RequestFacebookConversationID(senderID string, pageID string) (string, error) {
	access_token := os.Getenv("ACCESS_TOKEN")
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
