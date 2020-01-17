package utils

import "os"

// FileExists check file path for existent
func FileExists(fp string) bool {
	_, err := os.Stat(fp)
	return err == nil
}
