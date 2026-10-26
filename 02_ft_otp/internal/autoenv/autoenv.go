package autoenv

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type EnvEntry struct {
	Key   string
	Value string
}

// loads the .env file found at the root of the working directory and export it as enviroment variables
func parseLine(line string) (EnvEntry, error) {

	eq := strings.IndexRune(line, '=')
	if eq == -1 {
		return EnvEntry{}, fmt.Errorf("Missing assignment operator.")
	}
	// extract data
	key := line[:eq]
	value := line[eq+1:]

	// parse data
	value = strings.Trim(value, " 	\n\t\"'")
	key = strings.Trim(key, " 	\n\t")

	return EnvEntry{Key: key, Value: value}, nil
}

func parseEnvFile(envfile string) ([]EnvEntry, error) {

	file, err := os.Open(envfile)
	if err != nil {
		return nil, fmt.Errorf("Couldn't open .env file")
	}

	buf, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	entries := []EnvEntry{}

	split := strings.SplitSeq(string(buf), "\n")
	for line := range split {

		// trim spaces
		line = strings.Trim(line, " 	\n\t")
		// check comments
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		// parse
		entry, err := parseLine(line)
		if err != nil {
			return nil, err
		}
		// append
		entries = append(entries, entry)
	}

	return entries, nil
}

/* The only thing that will be run */
func Load(envfile string) error {
	entries, err := parseEnvFile(envfile)
	if err != nil {
		return err
	}

	for _, e := range entries {
		// fmt.Printf("Exporting '%s' -> '%s'\n", e.Key, e.Value)
		err := os.Setenv(e.Key, e.Value)
		if err != nil {
			return err
		}
	}
	return nil
}
