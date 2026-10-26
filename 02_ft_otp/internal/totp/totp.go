package totp

import (
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"time"

	"wasted-kiwi.com/ft_otp/internal/aes"
	"wasted-kiwi.com/ft_otp/internal/algo"
)

// The key argument should be the AES key, either 16 or 32 bytes
// to select AES-128 or AES-256.
var AES_KEY string = ""
var TIME_STEP uint64 = 30 // seconds

func LoadEnv() error {

	// AES_KEY
	aesKey := os.Getenv("AES_KEY")
	if aesKey == "" {
		return fmt.Errorf("No AES-GCM 32bytes encription key. Export the AES_KEY shell env variable.")
	}
	if len(aesKey) != 16 && len(aesKey) != 32 {
		return fmt.Errorf("Inavalid AES_KEY length (16 for AES-128, 32 for AES-256).")
	}
	AES_KEY = aesKey

	// TIME_STEP
	timeStep := os.Getenv("TIME_STEP")
	if timeStep == "" {
		// fmt.Printf("warn: No TIME_STEP found. Using default (30s).\n")
	} else {
		ts, err := strconv.Atoi(timeStep)
		if err != nil {
			return err
		}
		if ts == 0 {
			return fmt.Errorf("Invalid TIME_STP of 0.")
		}
		TIME_STEP = uint64(ts)
	}
	return nil
}

/*
 * Produce and encrypted ft_otp.key file containing the key (at least 64 bytes)
 * ft_otp.key: [N: chiphertext] + [16: Auth] + [12: nonce] + [16: EPOCH] = N + 44 bytes, where N is len(plaintext)
 */
func GenerateKey(secret []byte, outfile string) error {

	raw, err := hex.DecodeString(string(secret))

	// check if is 64 length hexadecimal
	if err != nil || len(raw) != 32 {
		return fmt.Errorf("key must be 64 hexadecimal characters.")
	}

	// cipher the decoded secret
	ciphertext, err := aes.Encrypt(raw, []byte(AES_KEY))
	if err != nil {
		return err
	}

	// write output
	err = os.WriteFile(outfile, ciphertext, 0600)
	if err != nil {
		return err
	}

	return nil
}

/*
 * Produces an OTP basesd on passed keyfile.
 * The time step is previously decided
 */
func ProduceOTP(path string) ([]uint32, error) {
	// open file
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// extract data
	// key := OTPKey{}
	// key.Parse(buf)

	// decipther text
	plaintext, err := aes.Decrypt(buf, []byte(AES_KEY))
	if err != nil {
		return nil, err
	}

	values := []uint32{}
	// generate OTP
	currStep := uint64(time.Now().UTC().Unix()) / TIME_STEP

	// Before
	values = append(values, algo.HOTP([]byte(plaintext), currStep-1))
	// Now
	values = append(values, algo.HOTP([]byte(plaintext), currStep))
	// Now +1
	values = append(values, algo.HOTP([]byte(plaintext), currStep+1))

	return values, nil
}

func StringOTP(otp uint32) string {
	return algo.StringOTP(otp)
}
