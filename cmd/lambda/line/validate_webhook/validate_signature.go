package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
)

var errInvalidSignature = errors.New("invalid signature")

func validateSignature(lineChannelSecret string, signature string, body string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("validateSignature: %w", err)
		}
	}()
	decoded, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return err
	}
	hash := hmac.New(sha256.New, []byte(lineChannelSecret))
	_, err = hash.Write([]byte(body))
	if err != nil {
		return err
	}
	if !hmac.Equal(decoded, hash.Sum(nil)) {
		return errInvalidSignature
	}
	return nil
}
