package cmd

import (
	"os"
)

// IsPathExists returns true if a given path exists, false if it doesn't.
// It might return an error if e.g. file exists but you don't have
// access.
func isPathExists(path string) (bool, error) {
	_, err := os.Lstat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	// error other than not existing e.g. permission denied
	return false, err
}