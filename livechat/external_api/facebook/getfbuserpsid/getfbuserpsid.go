package getfbuserpsid

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ParticipantsResponse struct {
	Participants Participants `json:"participants"`
}
type Participants struct {
	Data []Participant `json:"data"`
}

type Participant struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Id    string `json:"id"`
}

func GetUserPSID(accessToken string, pageID string, conversationID string) (_ string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("getfbuserprofile.GetUserPSID: %w", err)
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
