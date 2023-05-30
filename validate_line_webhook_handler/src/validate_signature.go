package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

func validateSignature(channelSecret string, signature string, body string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("validateSignature: %w", err)
		}
	}()
	decoded, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return err
	}
	hash := hmac.New(sha256.New, []byte(channelSecret))
	if _, err = hash.Write([]byte(body)); err != nil {
		return err
	}
	if !hmac.Equal(decoded, hash.Sum(nil)) {
		return err
	}
	return nil
}
