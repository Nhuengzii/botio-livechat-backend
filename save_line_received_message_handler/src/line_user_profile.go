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

func getLineUserProfile(userID string) (*lineUserProfile, error) {
	url := "https://api.line.me/v2/bot/profile/" + userID
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, &getLineUserProfileError{
			message: "couldn't get line user profile",
			err:     err,
		}
	}
	req.Header.Set("Authorization", "Bearer "+lineChannelAccessToken)
	resp, err := client.Do(req)
	if err != nil {
		return nil, &getLineUserProfileError{
			message: "couldn't get line user profile",
			err:     err,
		}
	}
	defer resp.Body.Close()
	usr := &lineUserProfile{}
	if err := json.NewDecoder(resp.Body).Decode(usr); err != nil {
		return nil, &getLineUserProfileError{
			message: "couldn't get line user profile",
			err:     err,
		}
	}
	if resp.StatusCode != http.StatusOK {
		return nil, &getLineUserProfileError{
			message: "couldn't get line user profile",
			err:     errors.New(usr.Message),
		}
	}
	return usr, nil
}

type getLineUserProfileError struct {
	message string
	err     error
}

func (e *getLineUserProfileError) Error() string {
	return e.message + ": " + e.err.Error()
}
