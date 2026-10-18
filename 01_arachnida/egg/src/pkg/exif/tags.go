package exif

import (
	"encoding/json"
	"fmt"
	"os"
)

// Define structs matching the JSON layout including values map
type ExifTag struct {
	ID          int               `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Values      map[string]string `json:"values"`
}

type ExifDatabase struct {
	Tags    map[string]ExifTag `json:"tags"`
	PTRTags map[string]ExifTag `json:"pointer_tags"`
	GPSTags map[string]ExifTag `json:"gps_tags"`
}

type TagLookup struct {
	String func(string) string
	Uint   func(uint16) string
	Value  func(uint16, any) string
}

var ExifLookup TagLookup

func init() {
	db, err := loadExifDB("./exif_tags.json")
	if err != nil {
		panic(err)
	}

	ExifLookup.String = getStringLookup(&db)
	ExifLookup.Uint = getUintLookup(&db)
	ExifLookup.Value = getValueLookup(&db)
}

func loadExifDB(path string) (ExifDatabase, error) {
	var db ExifDatabase

	data, err := os.ReadFile(path)
	if err != nil {
		return db, err
	}

	err = json.Unmarshal(data, &db)
	if err != nil {
		return db, err
	}

	return db, nil
}

func stringLookup(db ExifDatabase, hex string) string {
	if tag, exists := db.Tags[hex]; exists {
		return tag.Name
	} else if tag, exists := db.GPSTags[hex]; exists {
		return tag.Name
	} else if tag, exists := db.PTRTags[hex]; exists {
		return tag.Name
	} else {
		return fmt.Sprintf("UnknownTag(%s)", hex)
	}
}

func findTag(db ExifDatabase, hex string) (ExifTag, bool) {
	if tag, exists := db.Tags[hex]; exists {
		return tag, true
	}
	if tag, exists := db.GPSTags[hex]; exists {
		return tag, true
	}
	if tag, exists := db.PTRTags[hex]; exists {
		return tag, true
	}
	return ExifTag{}, false
}

func getStringLookup(db *ExifDatabase) func(string) string {
	return func(targetTagHex string) string {
		return stringLookup(*db, targetTagHex)
	}
}

func getUintLookup(db *ExifDatabase) func(uint16) string {
	return func(targetTagUint uint16) string {
		targetTagHex := fmt.Sprintf("0x%04X", targetTagUint)
		return stringLookup(*db, targetTagHex)
	}
}

func getValueLookup(db *ExifDatabase) func(uint16, any) string {
	return func(targetTagUint uint16, rawVal any) string {
		targetTagHex := fmt.Sprintf("0x%04X", targetTagUint)
		tag, exists := findTag(*db, targetTagHex)

		// Normalize rawVal into a slice of any
		var values []any
		switch v := rawVal.(type) {
		case []any:
			values = v
		case []byte:
			for _, b := range v {
				values = append(values, b)
			}
		default:
			values = []any{rawVal}
		}

		// If it's a single value, check if we have a readable lookup mapping for it
		if len(values) == 1 {
			valKey := fmt.Sprintf("%v", values[0])
			if exists && tag.Values != nil {
				if readable, ok := tag.Values[valKey]; ok {
					return readable
				}
			}
			return valKey
		}

		// If it contains multiple values (e.g. GPS coordinates, arrays), format them cleanly
		var formatted []string
		for _, val := range values {
			formatted = append(formatted, fmt.Sprintf("%v", val))
		}

		// You can join them or let fmt print the slice format
		return fmt.Sprintf("%v", formatted)
	}
}
