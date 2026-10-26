package main

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/hex"
	"net/url"
)

// otpauth://totp/Issuer:Account?secret=Secret&issuer=Issuer
const ISSUER = "topiana-totp"

func OTPAuthURL(hexKey, account string) (string, error) {
	raw, err := hex.DecodeString(hexKey)
	if err != nil {
		return "", err
	}

	secret := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(raw)

	label := url.PathEscape(ISSUER + ":" + account)

	q := url.Values{}
	q.Set("secret", secret)
	q.Set("issuer", ISSUER)
	// optional, these are the defaults:
	// q.Set("algorithm", "SHA1")
	// q.Set("digits", "6")
	// q.Set("period", "30")

	return "otpauth://totp/" + label + "?" + q.Encode(), nil
}

// 64 hexkey: 32 bytes -> hex encoded 64 bytes
func GenerateHexKey(length int) (string, error) {
	b := make([]byte, length/2)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
