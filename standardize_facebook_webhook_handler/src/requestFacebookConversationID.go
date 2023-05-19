package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func RequestFacebookConversationID(messageData MessageData, pageID string) (string, error) {
	// TODO: access_token still hardcode
	access_token := "EAACgkFuKQwcBAIWWdLrLpOYGJrrI2ZAQWfxolrzTjPFuxjZCOLMxXX8vH6rUhLs6sGB5X7aUBKLiBFzsoeBC13U8GpZAczfBosZBRYlSGigKAbYkzhAt46m8kpQAYoe3yWVSmnAl0xekyZC7Iw09eWM2XjJPKpW6PIhPBBFJh5Oz3tYxxSqe8"
	uri := fmt.Sprintf("https://graph.facebook.com/v16.0/%v/conversations?platform=Messenger&user_id=%v&access_token=%v",
		pageID, messageData.Sender.ID, access_token)
	discordLog("RequestFacebookConversationID")

	startTime := time.Now()
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

	discordLog(fmt.Sprint(body.Data))
	discordLog(fmt.Sprintf("ConversationID request elapsed : %v", time.Since(startTime)))
	return body.Data[0].Id, nil
}

type ResponseConversationID struct {
	Data []Conversation `json:"data"`
}

type Conversation struct {
	Id string `json:"id"`
}
