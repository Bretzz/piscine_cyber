package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"slices"
	"spider/pkg/fetch"
	"spider/pkg/robotstxt"
	"spider/pkg/scrape"
	"strings"
)

// SOURCE: https://gist.github.com/r0l1/3dcbb0c8f6cfe9c66ab8008f55f8f28b
// askForConfirmation asks the user for confirmation. A user must type in "yes" or "no" and
// then press enter. It has fuzzy matching, so "y", "Y", "yes", "YES", and "Yes" all count as
// confirmations. If the input is not recognized, it will ask again. The function does not return
// until it gets a valid response from the user.
func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		switch response {
		case "y", "yes":
			return true
		case "n", "no":
			return false
		}
	}
}

func download_image(dwPath string, rawurl string, r *http.Response, body []byte, verbose bool) error {

	dwurl, err := url.Parse(rawurl)
	if err != nil {
		return err
	}

	// 1. Build directory path
	dir := dwPath + dwurl.Host + filepath.Dir(dwurl.Path)
	dir = path.Clean(dir)

	// 2. Mkdir if not present
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	// 3. Get filename
	u, _ := url.Parse(r.Request.URL.String())
	filename := path.Base(u.Path)

	// 4. Create file
	file, err := os.Create(dir + "/" + filename)
	if err != nil {
		return err
	}

	// 5. Write file
	n, err := file.Write([]byte(body))
	if err != nil || n != len(body) {
		return err
	}

	if verbose == true {
		fmt.Printf("Image downloaded! %s\n", dir+"/"+filename)
	}
	return nil
}

type SpiderInfo struct {
	Host       string
	Robots     map[string]robotstxt.RobotsTXT
	Visited    []string
	Downloaded []string
}

// NewSpiderInfo automates the initialization
func NewSpiderInfo(host string) SpiderInfo {
	return SpiderInfo{
		Host:       host,
		Robots:     make(map[string]robotstxt.RobotsTXT),
		Visited:    []string{},
		Downloaded: []string{},
	}
}

func spider(_url *url.URL, flags SpiderFlags) {

	info := NewSpiderInfo(_url.Host)

	spiderling(_url, 0, &info, &flags)

	if flags.Verbose == true {
		// Print visited
		fmt.Println("-----------------------------------")
		fmt.Println("Visited:")
		for _, v := range info.Visited {
			fmt.Println(v)
		}
		// Print Downloaded
		fmt.Println("-----------------------------------")
		fmt.Println("Downloaded:")
		for _, v := range info.Downloaded {
			fmt.Println(v)
		}
	}

	// Print stats
	fmt.Println("-----------------------------------")
	fmt.Printf("%-11s %4d URLs\n%-11s %4d Images\n", "Visited", len(info.Visited), "Downloaded", len(info.Downloaded))

}

func spiderling(_url *url.URL, currLevel int, info *SpiderInfo, flags *SpiderFlags) {

	// #debug
	if flags.Verbose == true {
		fmt.Printf("Spidering URL '%s'\n", _url.String())
	}

	/* ---------------------------------------------- */
	/* 					ROBOTS.TXT					  */
	/* ---------------------------------------------- */
	// 1. Check if we visited the host, if not check host/robots.txt
	r, exist := info.Robots[_url.Host]
	// Add robot
	if exist == false {
		rob, err := robotstxt.Parse(_url)
		if err != nil {
			fmt.Println("[Warn] Parsing /robots.txt: ", err.Error())
			return
		}
		r = rob
		info.Robots[_url.Host] = r
	}
	// Check if route is allowed
	ret, err := r.IsAllowed(_url.Host, "*", _url.String())
	if err != nil {
		fmt.Printf("[Warn] Error validating /robots.txt Allowance for %s: %s\n", _url.String(), err.Error())
		return
	}
	if ret == false {
		if flags.Verbose == true {
			fmt.Printf("[Warn] Route blocked by /robots.txt: %s\n", _url.String())
		}
		return
	}

	/* ---------------------------------------------- */
	/* 				  FETCH & DOWNLOAD				  */
	/* ---------------------------------------------- */
	// 2. Fetch the URL
	body, resp, err := fetch.Raw(_url.String())
	if err != nil {
		fmt.Println("[Warn] Error while fetching the URL: ", err.Error())
		return
	}

	// Parse Content type
	ContentType := resp.Header.Get("Content-Type")
	if ContentType != "" {
		// needed for 'text/html; charset=utf-8' parsing
		ContentType = strings.Split(ContentType, ";")[0]
	}

	// 3. Check if Image
	if slices.Contains(targetTypes, ContentType) {
		err := download_image(flags.DownloadPath+"/"+info.Host+"/", _url.String(), resp, body, flags.Verbose)
		if err != nil {
			fmt.Println("[Warn] Error while writing file:", err.Error())
		} else {
			info.Downloaded = append(info.Downloaded, filepath.Base(_url.Path))
		}
		return
	}

	// save visited url
	info.Visited = append(info.Visited, _url.String())

	/* ---------------------------------------------- */
	/* 				  SCRAPE FURTHER				  */
	/* ---------------------------------------------- */
	// Break recursion
	if currLevel > flags.Limit {
		return
	}
	// Check if content is under Entripoint's host
	if flags.Spread == false && _url.Host != info.Host {
		return
	}
	// Check if contentType is scrapable
	if slices.Contains(scrapeTypes, ContentType) == false {
		return
	}

	// 4. find urls
	//------------------
	urls, err := scrape.HTML([]byte(body), "url")
	if err != nil {
		fmt.Println("[Warn] Error while reading the body")
		return
	}
	if urls == nil {
		// log.Fatal("No URL found")
		return
	}

	// 5. process found urls
	for _, foundUrl := range urls {
		// #debug
		// fmt.Printf("[%s] Found URL '%s'\n", prevUrl, url)

		// attach hostname if internal url
		if foundUrl != "" && strings.HasPrefix(foundUrl, "/") {
			foundUrl = "https://" + _url.Hostname() + foundUrl
		}

		// check if visited
		if slices.Contains(info.Visited, foundUrl) {
			continue
		}

		newUrl, err := url.Parse(foundUrl)
		if err != nil {
			fmt.Println("[Warn] Error while parsing URL: ", err.Error())
		} else {
			// spider the url
			spiderling(newUrl, currLevel+1, info, flags)
		}

	}
}

// const URL = "https://sdev.wasted-kiwi.com/"

// Response types to download
var targetTypes = []string{"image/jpg", "image/jpeg", "image/png", "image/gif", "image/bmp"}
var scrapeTypes = []string{"text/html", "application/javascript", "text/javascript"}

type SpiderFlags struct {
	Recursive    bool
	Limit        int
	DownloadPath string
	Verbose      bool
	Spread       bool
}

func main() {

	spidey := SpiderFlags{}

	// Parse flags
	recursive := flag.Bool("r", false, "Scrape the websites recursively")
	limit := flag.Int("l", 5, "Recursion depth to traverse")
	dwPath := flag.String("p", "./data/", "Path to store the images")
	verbose := flag.Bool("v", false, "Log the URLs while visiting")
	spread := flag.Bool("spread", false, "Navigate each url found")

	flag.Parse()

	if *recursive == false {
		*limit = 0
	}

	// Check ARG
	if flag.Arg(0) == "" {
		log.Fatal("Missing URL argument")
	}
	if len(flag.Args()) != 1 {
		log.Fatal("Too Many Arguments")
	}

	// Parse URL
	_url, err := url.Parse(flag.Arg(0))
	if err != nil {
		log.Fatal("Error while parsing URL: ", err.Error())
	}

	// Assign filags
	spidey.Recursive = *recursive
	spidey.Limit = *limit
	spidey.DownloadPath = *dwPath
	spidey.Verbose = *verbose
	spidey.Spread = *spread

	// #prod
	if !askForConfirmation(fmt.Sprintf("About to download images from URL '%s' into directory '%s'\nr_depth: %d, spread: %v\nProceed?", _url.String(), spidey.DownloadPath, spidey.Limit, spidey.Spread)) {
		return
	}

	// fmt.Printf("#DEBUG\n\tspider(%s, %s, %d)\n", *dwPath, _url.String(), *limit)
	spider(_url, spidey)
}
