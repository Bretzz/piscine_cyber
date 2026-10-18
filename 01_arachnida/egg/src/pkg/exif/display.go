package exif

import (
	"fmt"
)

func (e IFDEntry) toString() string {
	name := ExifLookup.Uint(e.tag)

	// Use ExifLookup.Value to get the human-readable string for this tag and raw value
	value := ExifLookup.Value(e.tag, e.pvalue)

	// Since ExifLookup.Value always returns a string, we can use %s here
	return fmt.Sprintf("%40s: [%s] %s\n", name, e.dataType.Name, value)
}

// Display all the info of the last parsed file
func Display() {
	if !parsed {
		fmt.Printf("No file parsed yet\n")
		return
	}

	for _, folder := range ifds {
		for _, entry := range folder.entries {
			fmt.Print(entry.toString())
		}
	}
}
