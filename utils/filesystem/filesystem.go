package filesystem

import (
	"os"
	"path/filepath"

	"github.com/lastrust/utils-go/logging"
)

// FileExists check file path for existent
func FileExists(fp string) bool {
	_, err := os.Stat(fp)
	return err == nil
}

func Remove(fp string) error {
	if err := os.Remove(fp); err != nil {
		logging.Err().Errorf("Failed to remove a file: %s %v\n", fp, err)
		return err
	} else {
		logging.Out().Infof("Removed a file: %s\n", fp)
		return nil
	}
}

func RemoveAll(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		logging.Err().Errorf("Failed to remove all files in a directory: %s %v\n", dir, err)
		return err
	} else {
		logging.Out().Infof("Removed all files in a directory: %s\n", dir)
		return nil
	}
}

func Move(old, new string) error {
	if err := os.Rename(old, new); err != nil {
		logging.Err().Errorf("Failed to remove a file from: %s to: %s %v\n", old, new, err)
		return err
	} else {
		logging.Out().Infof("Moved a file from: %s to: %s\n", old, new)
		return nil
	}
}

type File struct {
	Path string
	Info os.FileInfo
}

// GetFiles returns a files by directory path
func GetFiles(dir string) ([]File, error) {
	files := make([]File, 0)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, File{path, info})
		}
		return nil
	})
	return files, err
}

func TrimExt(filename string) string {
	return filename[0 : len(filename)-len(filepath.Ext(filename))]
}
