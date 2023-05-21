package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func RequestFacebookUserProfile(psid string) (ResponseFacebookUserProfile, error) {
	access_token := "EAACgkFuKQwcBAIWWdLrLpOYGJrrI2ZAQWfxolrzTjPFuxjZCOLMxXX8vH6rUhLs6sGB5X7aUBKLiBFzsoeBC13U8GpZAczfBosZBRYlSGigKAbYkzhAt46m8kpQAYoe3yWVSmnAl0xekyZC7Iw09eWM2XjJPKpW6PIhPBBFJh5Oz3tYxxSqe8"
	uri := fmt.Sprintf("https://graph.facebook.com/%v?fields=first_name,last_name,profile_pic&access_token=%v", psid, access_token)
	discordLog("RequestFacebookUserProfile")

	startTime := time.Now()
	resp, err := http.Get(uri)
	if err != nil {
		return ResponseFacebookUserProfile{}, err
	}
	defer resp.Body.Close()

	var body ResponseFacebookUserProfile
	err = json.NewDecoder(resp.Body).Decode(&body)

	discordLog(fmt.Sprintf("%+v", body))
	discordLog(fmt.Sprintf("UserProfileRequest Elasped : %v", time.Since(startTime)))
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
