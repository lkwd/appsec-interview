package fileops

import (
	"errors"
	"fmt"
	"os/exec"
	"path"
	"strings"
)

func checkPath(f string, expectedDir string) (string, error) {
	s := path.Clean(f)
	if strings.Contains(s, "/") && !strings.HasPrefix(s, expectedDir) {
		return "", errors.New("invalid path")
	}
	return s, nil
}

func Rename(parentDir string, file string, newName string) ([]byte, error) {

	checkedFile, err := checkPath(file, parentDir)
	if err != nil {
		return nil, err
	}

	checkedNewPath, err := checkPath(newName, parentDir)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("bash", "-c", fmt.Sprintf("mv %s %s", checkedFile, checkedNewPath))
	cmd.Dir = parentDir

	return cmd.Output()
}
