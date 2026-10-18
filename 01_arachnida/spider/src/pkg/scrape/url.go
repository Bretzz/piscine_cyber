package scrape

import (
	"net/url"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

const URL_CHARS string = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-._~:/?#[]@!$&'()*+,;="
const ALPHABETH string = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"
const ALPHA_NUM string = "0123456789qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"

// URL CHECKS
// ------------
// checks the url, returns false on error
type CheckFunc func(string) bool

func notShort(url string) bool {
	if len(url) <= 1 {
		return false
	}
	return true
}

func notNumber(url string) bool {
	if _, err := strconv.Atoi(url); err == nil {
		return false
	}
	return true
}

func hasLetters(url string) bool {
	if strings.ContainsAny(url, ALPHABETH) {
		return true
	}
	return false
}

func anumStart(url string) bool {
	if strings.ContainsRune(ALPHA_NUM, rune(url[1])) {
		return true
	}
	return false
}

/*
 * This functions finds both External and Internal URLs found in the response body.
 *
 * External URLs start with 'http://' or 'https://' and are easily spotted.
 * Internal URLs start with '/' and are harder to pinpoint.
 * To do so we implemented a modulable way to add checks on the pseudo-strings.
 * The check functions have this prototype
 * func (url string) bool
 * On 'false' return the url is invalid
 */

var (
	// External: scheme + host + path, stopping at anything that can't be in a URL
	// or that delimits it in HTML/JS/JSON.
	externalRegexp = regexp.MustCompile(`https?://[^\s"'\x60<>\\(){}\[\]|^]+`)
	// External: scheme (subsided) + host + path, eg: '//example.com/path/to/location'
	externalHiddenRegexp = regexp.MustCompile(`//[^\s"'\x60<>\\(){}\[\]|^]+`)

	// Internal: a quote/backtick, then "/" followed by a non-slash char.
	//  - excludes "//cdn.com" (protocol-relative)
	//  - excludes a lone "/" and things like "text/html" (no quote right before the slash)
	internalRegexp = regexp.MustCompile(`["'\x60](/[^/\s"'\x60<>\\][^\s"'\x60<>\\]*)`)
)

// Define the checks to make on internal urls strings
var URLchecks = []CheckFunc{notShort, notNumber, hasLetters, anumStart}

func find_url_in_html(body string) []string {

	var urls []string

	// EXTERNAL
	for _, m := range externalRegexp.FindAllString(body, -1) {
		m = strings.TrimRight(m, ".,;:!?")
		if i := strings.Index(m, "${"); i >= 0 { // template literal
			m = m[:i]
		}
		u, err := url.Parse(m)
		if err != nil || u.Host == "" {
			continue
		}
		// optional sanity check: real hosts have a dot (or are localhost)
		if !strings.Contains(u.Host, ".") && !strings.HasPrefix(u.Host, "localhost") {
			continue
		}
		urls = append(urls, m)
	}

	// EXTERNAL
	for _, m := range externalHiddenRegexp.FindAllString(body, -1) {
		m = m[2:] // //example.com -> example.com
		m = strings.TrimRight(m, ".,;:!?")
		if i := strings.Index(m, "${"); i >= 0 { // template literal
			m = m[:i]
		}

		u, err := url.Parse(m)
		if err != nil || u.Host == "" {
			continue
		}
		// optional sanity check: real hosts have a dot (or are localhost)
		if !strings.Contains(u.Host, ".") && !strings.HasPrefix(u.Host, "localhost") {
			continue
		}
		urls = append(urls, m)
	}

	// INTERNAL
	for _, m := range internalRegexp.FindAllStringSubmatch(body, -1) {
		p := m[1]
		if i := strings.Index(p, "${"); i >= 0 { // /api/users/${id} -> /api/users/
			p = p[:i]
		}
		if p == "" || p == "/" {
			continue
		}

		ok := true
		for _, f := range URLchecks {
			if !f(p) {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}

		/* ref, err := url.Parse(p)
		if err != nil {
			continue
		}
		urls = append(urls, base.ResolveReference(ref).String()) */
		urls = append(urls, p)
	}

	slices.Sort(urls)
	return slices.Compact(urls)
}
