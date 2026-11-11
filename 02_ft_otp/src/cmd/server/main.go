package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"wasted-kiwi.com/ft_otp/internal/totp"
)

// where to store keys
const KEY_DIR = "otpkeys"

func main() {
	// build key dir
	os.Mkdir(KEY_DIR, 0755)

	// Read env variables
	// load env variables
	{
		err := totp.LoadEnv()
		if err != nil {
			fmt.Printf("error: %s\n", err)
			os.Exit(1)
		}
	}

	mux := http.NewServeMux()

	// Backend APIS
	/*
	 * Request body: { account: "string" }
	 * Reply:
	 * 		201 -> text/plain (otpauth://...)
	 * 		400 -> text/plain (error in request)
	 *		502 -> text/plain (internal error)
	 */
	mux.HandleFunc("POST /api/otp/new", newOTP)
	/*
	 * Request body: { account: "string", otp: "6-digits-string" }
	 * Reply:
	 * 		200 -> OK
	 *		400 -> text/plain (error in request)
	 *		403 -> Invalid OTP
	 *		404 -> Account not found
	 *		502 -> text/plain (internal error)
	 */
	mux.HandleFunc("POST /api/otp/validate", validateOTP)

	// Server Frontend
	// --- Serve index.html on root GET / ---
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		// Optional: If someone visits a subpath that doesn't exist, you can return 404
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		// Serves the index.html file from the current directory
		http.ServeFile(w, r, "public/index.html")
	})

	fmt.Println("OTP server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
