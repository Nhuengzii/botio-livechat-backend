package main

import (
	"errors"
	"log"
	"strings"
)

const (
	headerNameXSign = "X-Hub-Signature"
	signaturePrefix = "sha1="
	appSecret       = "a55af6ed66089f33281b4d0963a2893b"
)

// errors
var (
	errNoXSignHeaders     = errors.New("There is no x-sign-header")
	errInvalidXSignHeader = errors.New("Invalid x-sign header")
)

func VertifySignature(header map[string]string) error {
	log.Println("Header : ", header)
	recieveSignature := header["x-hub-signature-256"]
	if !strings.HasPrefix(recieveSignature, signaturePrefix) {
		return errNoXSignHeaders
	}

	return nil
}
