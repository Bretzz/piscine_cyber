package find

import "fmt"

type Sequence []uint8

// Returns the byte offset to the start of the sequence found
func ByteSequence(b []byte, seq Sequence) (uint, error) {
	if len(seq) == 0 {
		return 0, fmt.Errorf("Zero length sequence")
	}

	// loop the bytes and the sequence
	for i := range b {
		ok := true
		for j := range seq {

			// Check if we finished the bytes
			if i+j == len(b) {
				return 0, fmt.Errorf("Sequence not found")
			}

			// Check if we are following the sequence
			if b[i+j] != byte(seq[j]) {
				ok = false
				break
			}
		}

		// check if the sequence was found
		if ok == true {
			return uint(i), nil
		}
	}

	return 0, fmt.Errorf("Sequence not found")
}

// expected 2 bytes. either II or MM
func IsLittle(b []byte) (bool, error) {

	if string(b) == "II" {
		return true, nil
	} else if string(b) == "MM" {
		return false, nil
	} else {
		return false, fmt.Errorf("Invalid endiannes specified '%s'. Expected 'II' or 'MM'\n", string(b))
	}
}
