package filesystem

import (
	"os"
	"path/filepath"
)

// FileExists check file path for existent
func FileExists(fp string) bool {
	_, err := os.Stat(fp)
	return err == nil
}

type File struct {
	Path string
	Info os.FileInfo
}

// GetFiles returns a files by directory path
func GetFiles(dir string) (files []File, err error) {
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, File{path, info})
		}
		return nil
	})
	return
}
