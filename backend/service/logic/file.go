package logic

import (
	"os"
	"path/filepath"
	"strings"
)

var Dir string
var ImageDir string

func InitFileSys() error {
	directory, err := os.Getwd()
	if err != nil {
		return err
	}
	Dir = directory

	directory = filepath.Join(directory, "dist", "images")
	if _, err := os.Stat(directory); err != nil {
		if derr := os.Mkdir(directory, os.ModePerm); derr != nil {
			return err
		}
	}
	ImageDir = directory
	return nil
}

func SaveImage(fileName string, fileBuf []byte) error {
	fullPath := filepath.Join(ImageDir, fileName)
	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(fileBuf)
	return err
}

func SaveImageForSub(subDir, fileName string, fileBuf []byte) error {
	fullPath := filepath.Join(ImageDir, fileName)
	if subDir != "" {
		fullPath = filepath.Join(ImageDir, subDir, fileName)
	}
	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(fileBuf)
	return err
}

func FindImage(tokenId string) string {
	entries, err := os.ReadDir(ImageDir)
	if err != nil {
		return tokenId + ".png"
	}

	for _, entry := range entries {
		if entry.Type().IsRegular() {
			if strings.HasPrefix(entry.Name(), tokenId) {
				return entry.Name()
			}
		}
	}
	return tokenId + ".png"
}
