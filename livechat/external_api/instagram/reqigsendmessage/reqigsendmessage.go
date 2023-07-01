// Package reqfbsendmessage implement a function to call graph API request for sending message.
//
// # Uses Graph API v16.0
package reqigsendmessage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SendMessage makes a graph API call to send message to target recipient and return SendingMessageResponse which is facebook response to the message sending request.
// Return an error if it occurs.
//
// # Allow sending one text message or one attachment message. Cannot be send together.
//
// *** note that sending picture or file or video might take some time ***
//
// Use instagram page accessToken.
func SendMessage(accessToken string, message SendingMessage, facebookPageID string) (_ *SendingMessageResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("reqigsendmessage.SendMessage: %w", err)
		}
	}()

	url := fmt.Sprintf("https://graph.facebook.com/v16.0/%v/messages?access_token=%v", facebookPageID, accessToken)
	instagramReqBody, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(instagramReqBody))
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

// SendingMessageResponse contains instagram's send message via api response
type SendingMessageResponse struct {
	RecipientID string `json:"recipient_id"` // IGSID of the recipient
	MessageID   string `json:"message_id"`   // sent message's ID
	Timestamp   int64  `json:"timestamp"`    // sent message's timestamp
}

// SendingMessage contains facebook send message request body
//
// [instagram Send API doc] : https://developers.facebook.com/docs/messenger-platform/instagram
type SendingMessage struct {
	Recipient Recipient `json:"recipient"` // target user IGSID
	Message   Message   `json:"message"`   // message body contains either text message or attachment message
}

// interface for TextMessage and Message with attachment
type Message interface {
	Message()
}

// MessageText contains text string that caller wanted to send
type MessageText struct {
	Text string `json:"text"` // text message
}

// MessageAttachment contains attachment that caller wanted to send
type MessageAttachment struct {
	Attachment AttachmentInstagramRequest `json:"attachment"` // attachment message
}

// implement Message
func (MessageText) Message()       {}
func (MessageAttachment) Message() {}

// AttachmentInstagramRequest contain informations about the request attachment
type AttachmentInstagramRequest struct {
	AttachmentType string                     `json:"type"`    // type of the attachment instagram supported
	Payload        AttachmentInstagramPayload `json:"payload"` // actual payload of the attachment
}

// AttachmentFacebookRequest contain informations about the attachment payload
type AttachmentInstagramPayload struct {
	Src          string     `json:"url,omitempty"`           // usable for media type attachment(image,video,audio,file). URL of the attachment.
	IsReusable   bool       `json:"is_reusable,omitempty"`   // the Messenger Platform supports saving assets via the Send API and Attachment Upload API. This allows you reuse assets, rather than uploading them every time they are needed.
	TemplateType string     `json:"template_type,omitempty"` // type of template caller want to send. only usable if the AttachmentInstagramRequest's AttachmentType is "template".
	Elements     []Template `json:"elements,omitempty"`      // each element must match the template type. only usable if the AttachmentInstagramRequest's AttachmentType is "template".
}

// Recipient contain target reciever of the message's informations
type Recipient struct {
	Id string `json:"id"` // IGSID of the reciever
}
