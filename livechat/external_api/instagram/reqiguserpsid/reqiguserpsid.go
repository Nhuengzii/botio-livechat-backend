// Package reqiguserpsid implement a function to make a graph api request for a user's psid in specific conversation.
//
// # Uses Graph API v16.0
package reqiguserpsid

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// A ParticipantsResponse is a response body recieved from graph API request which contains a Particiapants object
type ParticipantsResponse struct {
	Participants Participants `json:"participants"` // Participants of the conversation
}

// A Participants contains a slice of Participant object in a specific conversation
type Participants struct {
	Data []Participant `json:"data"` // data of the participant
}

// A Participant contains a partcipant's information
type Participant struct {
	Name  string `json:"name"`  // name of the participant
	Email string `json:"email"` // email of the participant
	Id    string `json:"id"`    // IGSID of the participant
}

// GetUserIGSID makes a graph API call and returns a string of specific conversation's participant user psid.
// Only return IGSID for a conversation in the specify page.
// Return an error if it occurs.
//
// Use instagram page accessToken.
func GetUserIGSID(accessToken string, pageID string, conversationID string) (_ string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("reqiguserprofile.GetUserPSID: %w", err)
		}
	}()

	url := fmt.Sprintf("https://graph.facebook.com/%v?fields=participants&access_token=%v", conversationID, accessToken)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var participantsResponse ParticipantsResponse
	err = json.NewDecoder(resp.Body).Decode(&participantsResponse)
	if err != nil {
		return "", err
	}

	var psid string
	for _, participant := range participantsResponse.Participants.Data {
		if participant.Id != pageID {
			psid = participant.Id
			break
		}
	}
	return psid, nil
}
