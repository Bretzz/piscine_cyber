package aes

// Credits: https://gist.github.com/kkirsche/e28da6754c39d5e7ea10

import (
	"io"

	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	// "encoding/hex"
)

// AES-GCM should be used because the operation is an authenticated encryption
// algorithm designed to provide both data authenticity (integrity) as well as
// confidentiality.

// Merged into Golang in https://go-review.googlesource.com/#/c/18803/

// The key argument should be the AES key, either 16 or 32 bytes
// to select AES-128 or AES-256.

type Reader struct{}

func (Reader) Read(p []byte) (n int, err error) {
	return rand.Read(p)
}

const NONCE_SIZE = 12

func Encrypt(plaintext, key []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	nonce := make([]byte, NONCE_SIZE)
	r := Reader{}
	if _, err := io.ReadFull(r, nonce); err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

	// fmt.Printf("plaintext: %d\n", len(plaintext))
	// fmt.Printf("ciphertext: %d\n", len(ciphertext))
	// fmt.Printf("nonce: %d\n", len(nonce))

	cipherblock := append(ciphertext, nonce...)

	// fmt.Printf("ciphertext: '%s'\n", ciphertext)
	// fmt.Printf("nonce: '%s'\n", nonce)
	// fmt.Printf("cipherblock: '%s'\n", cipherblock)

	return cipherblock, nil
}

func Decrypt(cipherblock, key []byte) ([]byte, error) {
	// The key argument should be the AES key, either 16 or 32 bytes
	// to select AES-128 or AES-256.

	// nonce, _ := hex.DecodeString("bb8ef84243d2ee95a41c6c57")

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := cipherblock[len(cipherblock)-NONCE_SIZE:]
	ciphertext := cipherblock[:len(cipherblock)-NONCE_SIZE]

	// fmt.Printf("cipherblock: '%s'\n", cipherblock)
	// fmt.Printf("ciphertext: '%s'\n", ciphertext)
	// fmt.Printf("nonce: '%s'\n", nonce)

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	// fmt.Printf("%s\n", string(plaintext))
	return plaintext, nil
}
