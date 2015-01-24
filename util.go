package main

import (
	"errors"
	"os"
	"path/filepath"
)

func openOrCreateFile(fileName string) (*os.File, error) {
	// TODO maybe should set the flag and FileMode by user.
	file, err := os.OpenFile(fileName, os.O_APPEND, 0644)
	if err != nil {
		dirPath := filepath.Dir(fileName)
		err = os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return nil, err
		}
		file, err = os.Create(fileName)
		if err != nil {
			return nil, errors.New("can not create the file.")
		}
	}
	return file, nil
}
