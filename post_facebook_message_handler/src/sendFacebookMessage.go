package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func SendFacebookMessage(requestMessage RequestMessage, psid string, response *FacebookResponse) error {
	token := os.Getenv("ACCESS_TOKEN")
	uri := fmt.Sprintf("https://graph.facebook.com/v16.0/me/messages?access_token=%v", token)

	facebookRequest := FacebookRequest{
		Recipient: Recipient{
			Id: psid,
		},
		Message: Message{
			Text: requestMessage.Message,
			Attachment: AttachmentFacebookRequest{
				AttachmentType: requestMessage.Attachment.AttachmentType,
				Payload: AttachmentFacebookPayload{
					Src:        requestMessage.Attachment.Payload.Src,
					IsReusable: true,
				},
			},
		},
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(facebookRequest)
	if err != nil {
		discordLog(fmt.Sprintf("Error encoding body : %v", err))
	}
	resp, err := http.Post(uri, "application/json", &buf)
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
	RecipientID string `json:"recipient_id"`
	MessageID   string `json:"message_id"`
}

type FacebookRequest struct {
	Recipient Recipient `json:"recipient"`
	Message   Message   `json:"message"`
}

type Message struct {
	Text       string                    `json:"text"`
	Attachment AttachmentFacebookRequest `json:"attachment"`
}
type AttachmentFacebookRequest struct {
	AttachmentType string                    `json:"type"`
	Payload        AttachmentFacebookPayload `json:"payload"`
}
type AttachmentFacebookPayload struct {
	Src        string `json:"url"`
	IsReusable bool   `json:"is_reusable"`
}
type Recipient struct {
	Id string `json:"id"`
}
