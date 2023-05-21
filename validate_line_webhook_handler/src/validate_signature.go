package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"os"
)

var channelSecret = os.Getenv("LINE_CHANNEL_SECRET")

func validateSignature(channelSecret string, signature string, body []byte) bool {
	decoded, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false
	}
	hash := hmac.New(sha256.New, []byte(channelSecret))
	_, err = hash.Write(body)
	if err != nil {
		return false
	}
	return hmac.Equal(decoded, hash.Sum(nil))
}
