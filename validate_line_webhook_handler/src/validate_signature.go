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
		return &validateSignatureError{
			message: "couldn't validate signature",
			err:     err,
		}
	}
	hash := hmac.New(sha256.New, []byte(channelSecret))
	_, err = hash.Write([]byte(body))
	if err != nil {
		return &validateSignatureError{
			message: "couldn't validate signature",
			err:     err,
		}
	}
	valid := hmac.Equal(decoded, hash.Sum(nil))
	if !valid {
		return &validateSignatureError{
			message: "couldn't validate signature",
			err:     errors.New("invalid signature"),
		}
	}
	return nil
}

type validateSignatureError struct {
	message string
	err     error
}

func (e *validateSignatureError) Error() string {
	return e.message + ": " + e.err.Error()
}
