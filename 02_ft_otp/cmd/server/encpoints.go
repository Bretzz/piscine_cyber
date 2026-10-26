package main

import (
	"encoding/json"
	"errors"
	"io/fs"
	"net/http"
	"path/filepath"
	"strconv"

	"wasted-kiwi.com/ft_otp/internal/totp"
)

// ft_otp interactions

func isdigit(s string) bool {
	if s == "" {
		return false
	}

	for i := range len(s) {
		c := s[i]
		isDigit := (c >= '0' && c <= '9')
		if !isDigit {
			return false
		}
	}

	return true
}

func isanum(s string) bool {
	if s == "" {
		return false
	}
	for i := range len(s) {
		c := s[i]
		isLetter := (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
		isDigit := (c >= '0' && c <= '9')
		if !isLetter && !isDigit {
			return false
		}
	}
	return true
}

func keypath(account string) string {
	return filepath.Join(KEY_DIR, account+".key")
}

// registry with keys, mutexes and validators
var registry totp.Registry

type RequestNewOTP struct {
	Account string `json:"account"`
}

func newOTP(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is correct
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Limit request body size to prevent memory attacks (e.g., 1MB)
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)

	// Decode the JSON request body
	var req RequestNewOTP
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	// parse account
	if isanum(req.Account) == false {
		http.Error(w, "Invalid account (only ASCII (a-z, A-Z, 0-9))", http.StatusBadRequest)
		return
	}

	// generate secret
	hexkey, err := GenerateHexKey(64)
	if err != nil {
		http.Error(w, "Oopsie, something is off...", http.StatusInternalServerError)
		return
	}

	// generate URL
	url, err := OTPAuthURL(hexkey, req.Account)
	if err != nil {
		http.Error(w, "Oopsie, something is off...", http.StatusInternalServerError)
		return
	}

	// #debug
	// fmt.Printf("URL: %s\n", url)

	// save .key
	err = totp.GenerateKey([]byte(hexkey), keypath(req.Account))
	if err != nil {
		http.Error(w, "Oopsie, something is off...", http.StatusInternalServerError)
		return
	}
	registry.Reset(req.Account, keypath(req.Account)) // new key => clear old cache/step

	// Success: keyfile created)
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(url))
}

type RequestValidateOTP struct {
	Account string `json:"account"`
	OTP     string `json:"otp"`
}

func validateOTP(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is correct
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Limit request body size to prevent memory attacks (e.g., 1MB)
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)

	// Decode the JSON request body
	var req RequestValidateOTP
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Check for errors
	if !isanum(req.Account) || len(req.OTP) != 6 || !isdigit(req.OTP) {
		http.Error(w, "Invali Parameters", http.StatusBadRequest)
		return
	}

	// check if key exist
	/* path := keypath(req.Account)
	_, err := os.Stat(path)
	if err != nil {
		// not found
		if errors.Is(err, fs.ErrNotExist) {
			http.Error(w, "Account not found", http.StatusNotFound)
			return
		}
	}

	// read stored otp
	checks, err := totp.ProduceOTP(path)
	if err != nil {
		http.Error(w, "Oopsie, something is off...", http.StatusInternalServerError)
		return
	}

	otp, err := strconv.Atoi(req.OTP)
	if err != nil {
		http.Error(w, "Oopsie, something is off...", http.StatusInternalServerError)
		return
	}

	// check stored against sent
	if slices.Contains(checks, uint32(otp)) {
		// Successful authentication
		_, _ = w.Write([]byte("Authenticated Successfully"))
		return
	}

	// Failed to auth
	http.Error(w, "Unauthorized", http.StatusUnauthorized) */

	otp, err := strconv.Atoi(req.OTP) // already validated as 6 digits
	if err != nil {
		http.Error(w, "Invalid Parameters", http.StatusBadRequest)
		return
	}

	ok, err := registry.Verify(req.Account, keypath(req.Account), uint32(otp))
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			http.Error(w, "Account not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Oopsie, something is off...", http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	_, _ = w.Write([]byte("Authenticated Successfully"))
}
