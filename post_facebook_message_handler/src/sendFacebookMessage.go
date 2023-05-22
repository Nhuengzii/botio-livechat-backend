package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func SendFacebookMessage(response *FacebookResponse) error {
	token := os.Getenv("ACCESS_TOKEN")
	uri := fmt.Sprintf("https://graph.facebook.com/v16.0/PAGE-ID/messages" +
		"?recipient={'id':'PSID'}" +
		"&messaging_type=RESPONSE" +
		"&message={'text':'hello,world'}" +
		"&access_token=PAGE-ACCESS-TOKEN")

	resp, err := http.Get(uri)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		return err
	}
	return nil
}

type FacebookResponse struct {
	RecipientID string `json:"recipient_id`
	MessageID   string `json:"message_id"`
}
