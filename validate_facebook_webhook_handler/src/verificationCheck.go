package main

import "errors"

func VerificationCheck(queryStringParameters map[string]string) error {
	mode := queryStringParameters["hub.mode"]
	token := queryStringParameters["hub.verify_token"]

	if config.verifyToken != token {
		return errors.New("verify_token : token does not match")
	} else if mode != "subscribe" {
		return errors.New("mode : mode is not subscibe")
	}

	return nil

}
