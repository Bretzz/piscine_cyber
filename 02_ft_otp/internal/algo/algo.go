package algo

import (
	"fmt"

	"crypto/hmac"
	"crypto/sha1"
	"encoding/binary"
	// "encoding/hex"
)

// Ref: https://www.rfc-editor.org/info/rfc4226/

// HOTP(K,C) = Truncate(HMAC-SHA-1(K,C))
/*
   Step 1: Generate an HMAC-SHA-1 value
   Let HS = HMAC-SHA-1(K,C)  // HS is a 20-byte string

   Step 2: Generate a 4-byte string (Dynamic Truncation)
   Let Sbits = DT(HS)   //  DT, defined below,
                        //  returns a 31-bit string

   Step 3: Compute an HOTP value
   Let Snum  = StToNum(Sbits)   // Convert S to a number in
                                    0...2^{31}-1
   Return D = Snum mod 10^Digit //  D is a number in the range
                                    0...10^{Digit}-1
*/
// Generates an Hash from the concatenated Key and Counter.
// The counter is rearranged as BigEndian
func private_HMAC_SHA_1(K []byte, C uint64) []byte {
	// convert to Big Endian
	be := make([]byte, 8)
	binary.BigEndian.PutUint64(be, C)

	// Initialize HMAC with the secret key (K)
	mac := hmac.New(sha1.New, K)

	// Write the counter (be) as the message payload
	mac.Write(be)

	// Compute and return the resulting hash bytes (20 bytes for SHA1)
	sum := mac.Sum(nil)

	// #debug
	// fmt.Println(hex.EncodeToString(sum))
	// fmt.Println(hs)
	// fmt.Println(len(sum), "bytes")

	return sum
}

// Dynamically Truncate the Hash string
/*
	int offset   =  hmac_result[19] & 0xf ;
	int bin_code = (hmac_result[offset]  & 0x7f) << 24
		| (hmac_result[offset+1] & 0xff) << 16
		| (hmac_result[offset+2] & 0xff) <<  8
		| (hmac_result[offset+3] & 0xff) ;
*/
func private_DT(HS []byte) uint32 {

	// Get the lower 4 bits of the last byte as the offset (0-15)
	offset := HS[19] & 0xf

	// Extract a 4-byte window starting at the offset and convert to uint32
	binCode := (uint32(HS[offset]&0x7f) << 24) |
		(uint32(HS[offset+1]&0xff) << 16) |
		(uint32(HS[offset+2]&0xff) << 8) |
		uint32(HS[offset+3]&0xff)

	return binCode
}

func HOTP(K []byte, C uint64) uint32 {

	hs := private_HMAC_SHA_1(K, C)
	// fmt.Println(hs)

	bc := private_DT(hs)
	// fmt.Println(bc)

	otp := bc % 1000000 // 10^6 for Digits = 6

	return otp
}

func StringOTP(otp uint32) string {
	return fmt.Sprintf("%06d", otp)
}
