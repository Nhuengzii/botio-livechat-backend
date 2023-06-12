package getfbuserprofile

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type UserProfile struct {
	Name       string `json:"name"`
	ProfilePic string `json:"profile_pic"`
	Locale     string `json:"locale"`
	TimeZone   string `json:"timezone"`
	Gender     string `json:"gender"`
}

func GetUserProfile(accessToken string, psid string) (_ *UserProfile, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("getfbuserprofile.GetUserProfile: %w", err)
		}
	}()
	url := fmt.Sprintf("https://graph.facebook.com/%v?fields=name,profile_pic&access_token=%v", psid, accessToken)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var userProfile *UserProfile
	err = json.NewDecoder(resp.Body).Decode(userProfile)
	if err != nil {
		return nil, err
	}
	return userProfile, nil
}
