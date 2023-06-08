package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
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
	appSecret             = os.Getenv("APP_SECRET")
)

func VerifySignature(header map[string]string, bodyByte []byte) error {
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
