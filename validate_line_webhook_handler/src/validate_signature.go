package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

func validateSignature(channelSecret string, signature string, body string) error {
	decoded, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return err
	}
	hash := hmac.New(sha256.New, []byte(channelSecret))
	_, err = hash.Write([]byte(body))
	if err != nil {
		return err
	}
	valid := hmac.Equal(decoded, hash.Sum(nil))
	if !valid {
		return errors.New("signature invalid")
	}
	return nil
}
