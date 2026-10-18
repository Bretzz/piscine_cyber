package robotstxt

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// Route : /api/v1
type RobotAgent struct {
	ID       string
	Allow    []string
	Disallow []string
}

type RobotsTXT struct {
	Host     string
	Agents   map[string]RobotAgent
	Sitemaps []string
}

// NewRobotsTXT automates the initialization
func NewRobotsTXT() RobotsTXT {
	return RobotsTXT{
		Agents: make(map[string]RobotAgent),
	}
}

// Find the agent (Rudimental)
func (r *RobotsTXT) find_agent(agentId string) RobotAgent {
	for key, value := range r.Agents {
		trim := strings.Trim(key, "*")
		if trim != "" && strings.HasPrefix(agentId, trim) {
			return value
		}
	}

	standard, exist := r.Agents["*"]
	if exist == true {
		return standard
	}

	return RobotAgent{}
}

// rawurl: 'https://example.com/api/v1/
func (r *RobotsTXT) IsAllowed(host string, agentId string, rawurl string) (bool, error) {

	if host != r.Host {
		return false, fmt.Errorf("Wrong host: %s", host)
	}

	// check if no agents (everything allowed)
	if len(r.Agents) == 0 {
		return true, nil
	}

	// Find corresponding agent
	agent := r.find_agent(agentId)

	// Check if url specifically allowed
	for _, allow := range agent.Allow {
		ret, err := url_is_under(rawurl, allow)
		if err != nil {
			return false, err
		}
		if ret == true {
			return true, nil
		}
	}

	// Check if url disallowed
	for _, disallow := range agent.Disallow {
		ret, err := url_is_under(rawurl, disallow)
		if err != nil {
			return false, err
		}
		if ret == true {
			return false, nil
		}
	}

	return true, nil
}

// IsUnder checks if the target URL is under the base URL.
func url_is_under(targetStr, baseStr string) (bool, error) {

	// #debug
	// fmt.Printf("Under '%s' == '%s'\n", targetStr, baseStr)

	target, err := url.Parse(targetStr)
	if err != nil {
		return false, err
	}

	base, err := url.Parse(baseStr)
	if err != nil {
		return false, err
	}

	// Scheme and Host must match exactly
	if target.Scheme != base.Scheme || target.Host != base.Host {
		return false, nil
	}

	// #debug
	// fmt.Printf("'%s' -- '%s'\n", base.Path, target.Path)

	// Normalize paths by trimming trailing slashes
	basePath := strings.TrimRight(base.Path, "/")
	targetPath := target.Path

	// If base path is root or empty
	if basePath == "" {
		return strings.HasPrefix(targetPath, "/"), nil
	}

	// Target path must equal base path or start with basePath + "/"
	return targetPath == basePath || strings.HasPrefix(targetPath, basePath+"/"), nil
}

/*
 * Parse the scheme://host/robots.txt file
 * Returns a RobotsTXT struct with all data
 */

func fetch_robots_txt(target *url.URL) (string, error) {

	robotsurl := "https://" + target.Host + "/robots.txt"

	resp, err := http.Get(robotsurl)
	if err != nil {
		return "<3", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "<3", nil
	}

	return string(body), nil
}

/*
Define parsing regexp
!USE THEM LINE BY LINE!
# todo
match 'Disallow: ' that allows everything
*/
var (
	/* Very forgiving regexp that allows comments (bash comments '#') */
	agentRegexp = regexp.MustCompile(`(?i)^[ \t]*User-agent[ \t]*:[ \t]*([^\s#]*)`)

	/* Also very forgiving regexp to match Disallow Routes */
	disallowRegexp = regexp.MustCompile(`(?i)^[ \t]*Disallow[ \t]*:[ \t]*([^\s#]*)`)

	/* Also very forgiving regexp to match Allow Routes */
	allowRegexp = regexp.MustCompile(`(?i)^[ \t]*Allow[ \t]*:[ \t]*([^\s#]*)`)

	/* Also very forgiving regexp to match Sitemaps */
	sitemapRegexp = regexp.MustCompile(`(?i)^[ \t]*Sitemap[ \t]*:[ \t]*(https?://[^\s"'\x60<>\\(){}\[\]|^]+)`)
)

func Parse(target *url.URL) (RobotsTXT, error) {
	// Init main struct
	robotstxt := NewRobotsTXT()
	robotstxt.Host = target.Host

	body, err := fetch_robots_txt(target)
	if err != nil {
		return robotstxt, err
	}

	agent := RobotAgent{}

	for i, line := range strings.Split(body, "\n") {

		// if we found an agent
		if a := agentRegexp.FindString(line); a != "" {

			// save previous agent
			if agent.ID != "" {
				// #prod multiple agents overwrite the first
				robotstxt.Agents[agent.ID] = agent
				agent = RobotAgent{}
			}

			// Initialize agent
			agent.ID = a
			continue
		}

		// check for allow route
		if allow := allowRegexp.FindString(line); allow != "" {
			if agent.ID == "" {
				fmt.Printf("[Warn] Allow route outside agent specification (line %d)\n", i)
				continue
			}
			agent.Allow = append(agent.Allow, allow)
			continue
		}

		// check for disallow route
		if disallow := disallowRegexp.FindString(line); disallow != "" {
			if agent.ID == "" {
				fmt.Printf("[Warn] Disallow route outside agent specification (line %d)\n", i)
				continue
			}
			agent.Disallow = append(agent.Disallow, disallow)
			continue
		}

		// check for sitemaps
		if sitemap := sitemapRegexp.FindString(line); sitemap != "" {
			robotstxt.Sitemaps = append(robotstxt.Sitemaps, sitemap)
			continue
		}
	}

	return robotstxt, nil
}
