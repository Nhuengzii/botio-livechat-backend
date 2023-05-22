package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func SendFacebookMessage(requestMessage RequestMessage, psid string, pageID string, response *FacebookResponse) error {
	token := os.Getenv("ACCESS_TOKEN")
	uri := fmt.Sprintf("https://graph.facebook.com/v16.0/%v/messages?access_token=%v", pageID, token)

	discordLog(fmt.Sprint(uri))
	var facebookRequest FacebookRequest
	if requestMessage.Message != "" {
		facebookRequest = FacebookRequest{
			Recipient: Recipient{
				Id: psid,
			},
			MessagingType: "RESPONSE",
			Message: MessageText{
				Text: requestMessage.Message,
			},
		}
	} else {
		facebookRequest = FacebookRequest{
			Recipient: Recipient{
				Id: psid,
			},
			MessagingType: "RESPONSE",
			Message: MessageAttachment{
				Attachment: AttachmentFacebookRequest{
					AttachmentType: requestMessage.Attachment.AttachmentType,
					Payload: AttachmentFacebookPayload{
						Src:        requestMessage.Attachment.Payload.Src,
						IsReusable: true,
					},
				},
			},
		}
	}
	discordLog(fmt.Sprintf("%+v", facebookRequest))
	facebookReqBody, err := json.Marshal(facebookRequest)
	if err != nil {
		discordLog(fmt.Sprintf("Error marshal body : %v", err))
	}

	resp, err := http.Post(uri, "application/json", bytes.NewReader(facebookReqBody))
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
	Recipient     Recipient `json:"recipient"`
	Message       any       `json:"message"`
	MessagingType string    `json:"messaging_type"`
}

type MessageText struct {
	Text string `json:"text"`
}

type MessageAttachment struct {
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
