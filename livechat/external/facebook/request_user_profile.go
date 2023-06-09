package fbrequest

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func RequestFacebookUserProfile(access_token string, psid string) (FacebookUserProfile, error) {
	uri := fmt.Sprintf("https://graph.facebook.com/%v?fields=name,profile_pic&access_token=%v", psid, access_token)

	resp, err := http.Get(uri)
	if err != nil {
		return FacebookUserProfile{}, err
	}
	defer resp.Body.Close()

	var body FacebookUserProfile
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return FacebookUserProfile{}, err
	}

	return body, nil
}

type FacebookUserProfile struct {
	Name       string `json:"name"`
	ProfilePic string `json:"profile_pic"`
	Locale     string `json:"locale"`
	TimeZone   string `json:"timezone"`
	Gender     string `json:"gender"`
}
