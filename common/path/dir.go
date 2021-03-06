// Package path contains utility functions for directory and path operations.
package path

import (
	"fmt"
	"io/ioutil"
	"os"
)

// Writable checks, if the passed path is writable It creates a temporary
// file within the path and tries to delete it afterwards.
func Writable(path string) (bool, error) {
	fp, err := ioutil.TempFile(path, "writable")
	if err != nil {
		return false, err
	}
	fname := fp.Name()
	fp.Close()
	return true, os.Remove(fname)
}

// CreateDir creates a directory at the given path.
func CreateDir(path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		// Path does not exist, create the directory
		err = os.Mkdir(path, os.ModeDir|os.FileMode(0755))
		return err
	}
	if !fi.IsDir() {
		return fmt.Errorf("'%s' is not a directory", path)
	}
	return nil
}

// DirExists checks, if the specified file path represents a directory.
func DirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return info.IsDir(), nil
}
