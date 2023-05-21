package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func RequestFacebookUserProfile(psid string) (ResponseFacebookUserProfile, error) {
	access_token := os.Getenv("ACCESS_TOKEN")
	uri := fmt.Sprintf("https://graph.facebook.com/%v?fields=first_name,last_name,profile_pic&access_token=%v", psid, access_token)

	resp, err := http.Get(uri)
	if err != nil {
		return ResponseFacebookUserProfile{}, err
	}
	defer resp.Body.Close()

	var body ResponseFacebookUserProfile
	err = json.NewDecoder(resp.Body).Decode(&body)

	return body, nil
}

type ResponseFacebookUserProfile struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	ProfilePic string `json:"profile_pic"`
	Locale     string `json:"locale"`
	TimeZone   string `json:"timezone"`
	Gender     string `json:"gender"`
}
