// Package reqfbuserprofile implement a function to call graph api request for a UserProfile information.
//
// # Uses Graph API v16.0
package reqfbuserprofile

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// A UserProfile contains information about specific facebook user.
type UserProfile struct {
	Name       string `json:"name"`        // name of the user. It is the firstname combines with the lastname
	ProfilePic string `json:"profile_pic"` // profile picture of the user contains a URL to that picture
	Locale     string `json:"locale"`      // locale zone of the user
	TimeZone   string `json:"timezone"`    // timezone of the user
	Gender     string `json:"gender"`      // gender of the user
}

// GetUserProfile makes a graph API call and returns a UserProfile of facebook specific facebook's user,If there is a user with matching PSID.
// Only return the user of a specific page.
// Return an error if it occurs.
//
// Use facebook page accessToken.
func GetUserProfile(accessToken string, psid string) (_ *UserProfile, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("reqfbuserprofile.GetUserProfile: %w", err)
		}
	}()
	url := fmt.Sprintf("https://graph.facebook.com/%v?fields=name,profile_pic&access_token=%v", psid, accessToken)
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
