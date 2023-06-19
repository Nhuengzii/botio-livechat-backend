package postfbmessage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func SendMessage(accessToken string, message SendingMessage, pageID string) (_ *SendingMessageResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("postfbmessage.SendMessage: %w", err)
		}
	}()
	url := fmt.Sprintf("https://graph.facebook.com/v16.0/%v/messages?access_token=%v", pageID, accessToken)
	facebookReqBody, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(facebookReqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	now := time.Now().UnixMilli()
	var response SendingMessageResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	response.Timestamp = now
	if err != nil {
		return nil, err
	}
	return &response, nil
}

type SendingMessageResponse struct {
	RecipientID string `json:"recipient_id"`
	MessageID   string `json:"message_id"`
	Timestamp   int64  `json:"timestamp"`
}

type SendingMessage struct {
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
	Src          string `json:"url,omitempty"`
	IsReusable   bool   `json:"is_reusable,omitempty"`
	TemplateType string `json:"template_type,omitempty"`
	Elements     []any  `json:"elements,omitempty"` // each element must match the template type
}

type Recipient struct {
	Id string `json:"id"`
}
