package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

func FileExists(fileName string) bool {
	_, err := os.Stat(fileName)
	return !errors.Is(err, os.ErrNotExist)
}

func Unwrap[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}

func readConfig(path string) (Sources, error) {
	configFile, err := os.Open(path)
	if err != nil {
		return Sources{}, err
	}
	defer configFile.Close()
	byteValue, _ := io.ReadAll(configFile)
	var sources Sources
	json.Unmarshal(byteValue, &sources)
	return sources, nil
}

func isNew(url string, oldArticles OldArticles) bool {
	for _, element := range oldArticles.Urls {
		if element == url {
			return false
		}
	}
	return true
}

func deprecateArticle(url string) error {
	oldArticles, err := readOldArticles()
	if err != nil {
		return err
	}
	var updatedOldArticles OldArticles
	updatedOldArticles.Urls = append(oldArticles.Urls, url)
	saveOldArticles(updatedOldArticles)
	return nil
}

func saveOldArticles(oldArticles OldArticles) error {
	jsonResult, err := json.Marshal(oldArticles)
	if err != nil {
		fmt.Println(err)
		return err
	}
	file, err := os.OpenFile(oldArticlesFileName, os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return err
	}
	file.Write(jsonResult)
	return nil
}
