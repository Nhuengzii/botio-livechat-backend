package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type userProfile struct {
	DisplayName   string `json:"displayName"`
	UserID        string `json:"userId"`
	PictureURL    string `json:"pictureUrl"`    // not included when user has no profile pic
	StatusMessage string `json:"statusMessage"` // not included when user doesn't have status message
	Message       string `json:"message"`       // only included in case of error
}

func getUserProfile(channelAccessToken string, userID string) (_ *userProfile, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("lambda/line/save_received_message/main.getUserProfile: %w", err)
		}
	}()
	url := "https://api.line.me/v2/bot/profile/" + userID
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+channelAccessToken)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var profile userProfile
	err = json.NewDecoder(resp.Body).Decode(&profile)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(profile.Message)
	}
	return &profile, nil
}
