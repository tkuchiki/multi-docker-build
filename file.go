package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func GoWalk(location string) chan string {
	ch := make(chan string)
	go func() {
		err := filepath.Walk(location, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() || filepath.Base(path) != "Dockerfile" {
				return nil
			}

			abspath, _ := filepath.Abs(path)

			ch <- abspath
			return nil
		})

		if err != nil {
			log.Fatal(err)
		}

		defer close(ch)
	}()
	return ch
}

func IsDir(path string) bool {
	fi, err := os.Stat(path)

	if err != nil {
		log.Fatal(err)
	}

	return fi.IsDir()
}

func FileNotExists(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}

func FileExists(path string) bool {
	return !FileNotExists(path)
}

func Symlink(destpath string) error {
	abspath, err := filepath.Abs(destpath)

	if err != nil {
		log.Fatal(err)
	}

	return os.Symlink(abspath, "Dockerfile")
}

func RemoveFile(path string) error {
	return os.Remove(path)
}

func IsCurrentDir(path string) bool {
	return path == "."
}

func GetImage(path string) string {
	return filepath.Base(filepath.Dir(path))
}

func DecodeJson(path string) interface{} {
	b, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	dec := json.NewDecoder(bytes.NewReader(b))

	var jsonData interface{}
	dec.Decode(&jsonData)

	return jsonData
}
