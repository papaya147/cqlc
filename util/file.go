package util

import (
	"fmt"
	"os"
	"strings"
)

func CheckPathExists(path string) error {
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("the path %s does not exist", path)
	}
	return nil
}

func GetFile(path string) (*os.File, error) {
	if err := CheckPathExists(path); err != nil {
		return nil, err
	}
	return os.Open(path)
}

func GetFilesInDir(path, extension string) ([]string, error) {
	if err := CheckPathExists(path); err != nil {
		return nil, err
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var fileNames []string
	for _, file := range files {
		fileName := file.Name()
		fileSplit := strings.Split(fileName, ".")
		if fileSplit[len(fileSplit)-1] != extension {
			continue
		}

		fileNames = append(fileNames, fileName)
	}

	return fileNames, nil
}

func GetFileContents(path string) ([]byte, error) {
	file, err := GetFile(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return os.ReadFile(path)
}

func WriteFile(path string, contents []byte) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(contents)
	return err
}

func DeleteFile(path string) error {
	return os.Remove(path)
}
