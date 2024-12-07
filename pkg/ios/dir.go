package ios

import (
	"fmt"
	"os"
	"path/filepath"
)

func MkdirAll(dir string, perm os.FileMode) (string, error) {
	var err error
	dir, err = filepath.Abs(dir)
	if err != nil {
		return "", fmt.Errorf(`invalid directory "%s"`, dir)
	}
	err = os.MkdirAll(dir, perm)
	return dir, err
}
