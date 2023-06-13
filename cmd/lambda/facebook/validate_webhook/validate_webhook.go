package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
)

const (
	headerNameXSign = "X-Hub-Signature"
	signaturePrefix = "sha256="
)

// errors
var (
	errNoXSignHeaders     = errors.New("there is no x-sign-header")
	errInvalidXSignHeader = errors.New("invalid x-sign header")
	errHexDecodeString    = errors.New("error decoding receiveSignature")
)

func VerifyMessageSignature(header map[string]string, bodyByte []byte, appSecret string) error {
	// use for facebook post request
	receiveSignature := header["X-Hub-Signature-256"]
	if !strings.HasPrefix(receiveSignature, signaturePrefix) {
		return errNoXSignHeaders
	}

	appSecretHmac := hmac.New(sha256.New, []byte(appSecret))
	_, err := appSecretHmac.Write(bodyByte)
	if err != nil {
		return err
	}
	expectedSignatureByte := appSecretHmac.Sum(nil)

	err = compareSignature(receiveSignature, expectedSignatureByte)
	if err != nil {
		return err
	}
	return nil
}

func compareSignature(receiveSignature string, expectedSignatureByte []byte) error {
	actualSignature, err := hex.DecodeString(strings.Split(receiveSignature, "=")[1])
	if err != nil {
		return errHexDecodeString
	}
	if !hmac.Equal(actualSignature, expectedSignatureByte) {
		return errInvalidXSignHeader
	}
	return nil
}

func VerifyConnection(queryStringParameters map[string]string, verificationString string) error {
	// use for facebook get request
	mode := queryStringParameters["hub.mode"]
	token := queryStringParameters["hub.verify_token"]

	if verificationString != token {
		return errors.New("verify_token : token does not match")
	} else if mode != "subscribe" {
		return errors.New("mode : mode is not subscribe")
	}

	return nil
}
