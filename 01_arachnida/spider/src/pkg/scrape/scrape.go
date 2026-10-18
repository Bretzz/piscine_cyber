package scrape

import (
	"fmt"
)

/* This package scrapes a file searching for various data */

func HTML(b []byte, target string) ([]string, error) {

	switch target {
	case "url":
		return find_url_in_html(string(b)), nil
	}

	return nil, fmt.Errorf("Invalid target '%s'\n", target)
}
