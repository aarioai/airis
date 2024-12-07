package ios

import (
	"os"
	"path/filepath"
)

func PrepareFile(filename string, dirPerm os.FileMode) (string, string, error) {
	filename, err := filepath.Abs(filename)
	if err != nil {
		return "", "", err
	}
	dir := filepath.Dir(filename)
	err = os.MkdirAll(dir, dirPerm)
	return filename, dir, err
}
