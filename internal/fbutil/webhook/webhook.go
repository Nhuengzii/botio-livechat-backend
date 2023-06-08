package webhook

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
	errNoXSignHeaders     = errors.New("There is no x-sign-header")
	errInvalidXSignHeader = errors.New("Invalid x-sign header")
	errHexDecodeString    = errors.New("Error decoding recieveSignature")
)

func VerifyMessageSignature(header map[string]string, bodyByte []byte, appSecret string) error {
	// use for facebook post request
	recieveSignature := header["X-Hub-Signature-256"]
	if !strings.HasPrefix(recieveSignature, signaturePrefix) {
		return errNoXSignHeaders
	}

	appSecretHmac := hmac.New(sha256.New, []byte(appSecret))
	_, err := appSecretHmac.Write(bodyByte)
	if err != nil {
		return err
	}
	expectedSignatureByte := appSecretHmac.Sum(nil)

	err = compareSignature(recieveSignature, expectedSignatureByte)
	if err != nil {
		return err
	}
	return nil
}

func compareSignature(recieveSignature string, expectedSignatureByte []byte) error {
	actualSignature, err := hex.DecodeString(strings.Split(recieveSignature, "=")[1])
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
		return errors.New("mode : mode is not subscibe")
	}

	return nil
}
