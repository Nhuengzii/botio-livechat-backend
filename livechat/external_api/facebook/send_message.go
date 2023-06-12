package facebook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func SendMessage(accessToken string, fbRequest FBSendMsgRequest, pageID string) (*FBSendMsgResponse, error) {
	uri := fmt.Sprintf("https://graph.facebook.com/v16.0/%v/messages?access_token=%v", pageID, accessToken)

	facebookReqBody, err := json.Marshal(fbRequest)
	if err != nil {
		return &FBSendMsgResponse{}, err
	}

	var response FBSendMsgResponse
	resp, err := http.Post(uri, "application/json", bytes.NewReader(facebookReqBody))
	if err != nil {
		return &FBSendMsgResponse{}, err
	}
	defer resp.Body.Close()

	now := time.Now().UnixMilli()
	err = json.NewDecoder(resp.Body).Decode(&response)

	response.Timestamp = now
	if err != nil {
		return &FBSendMsgResponse{}, err
	}
	return &response, nil
}

type FBSendMsgResponse struct {
	RecipientID string `json:"recipient_id"`
	MessageID   string `json:"message_id"`
	Timestamp   int64  `json:"timestamp"`
}

type FBSendMsgRequest struct {
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
