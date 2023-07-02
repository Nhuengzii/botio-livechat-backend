// Package reqiguserprofile implement a function to make a graph api request for a UserProfile information.
//
// # Uses Graph API v16.0
package reqiguserprofile

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// A UserProfile contains information about specific instagram user.
type UserProfile struct {
	Name       string `json:"name"`        // name of the user. It is the firstname combines with the lastname
	ProfilePic string `json:"profile_pic"` // profile picture of the user contains a URL to that picture
}

// GetUserProfile makes a graph API call and returns a UserProfile of instagram specific instagram's user,If there is a user with matching IGSID.
// Only return the user of a specific page.
// Return an error if it occurs.
//
// Use instagram page accessToken.
func GetUserProfile(accessToken string, psid string) (_ *UserProfile, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("reqiguserprofile.GetUserProfile: %w", err)
		}
	}()
	url := fmt.Sprintf("https://graph.facebook.com/%v?access_token=%v", psid, accessToken)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var userProfile UserProfile
	err = json.NewDecoder(resp.Body).Decode(&userProfile)
	if err != nil {
		return nil, err
	}
	return &userProfile, nil
}
