package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

type lineUserProfile struct {
	DisplayName   string `json:"displayName"`
	UserID        string `json:"userId"`
	PictureURL    string `json:"pictureUrl"`
	StatusMessage string `json:"statusMessage"`
	Message       string `json:"message"` // in case of error
}

func getLineUserProfile(userID string) (_ *lineUserProfile, err error) {
	defer func() {
		if err != nil {
			err = errors.New("getLineUserProfile: " + err.Error())
		}
	}()
	url := "https://api.line.me/v2/bot/profile/" + userID
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+lineChannelAccessToken)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	usr := &lineUserProfile{}
	if err := json.NewDecoder(resp.Body).Decode(usr); err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}
	return usr, nil
}
