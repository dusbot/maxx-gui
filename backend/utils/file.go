package utils

import (
	"fmt"
	"os"
	"path"
)

func CreateDirIfNotExists(dirPath string) error {
	_, err := os.Stat(dirPath)
	if err == nil {
		fmt.Printf("The dirPath[%s] exists", dirPath)
		return nil
	}

	if os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			return err
		}
		return nil
	}
	return err
}

func CreateDirUnderHomeIfNotExists(dirName string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dirPath := path.Join(homeDir, dirName)
	return dirPath, CreateDirIfNotExists(dirPath)
}

func FileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	if err != nil {
		fmt.Printf("os stat filepath[%s] with error[%+v]\n", filepath, err)
		return false
	}
	return true
}
