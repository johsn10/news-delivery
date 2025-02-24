package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/mmcdole/gofeed"
)

type Sources struct {
	Sources []Source `json:"sources"`
}

type Source struct {
	Url   string   `json:"url"`
	Prefs []string `json:"prefs"`
}

type OldArticles struct {
	Urls []string `json:"urls"`
}

type Article struct {
	Url         string
	Title       string
	Description string
	Content     string
}

var oldArticlesFileName = "oldArticles.json"

func main() {
	llm, err := initLLM()
	if err != nil {
		fmt.Println("error in initLMM\n", err)
		return
	}
	sources, err := readConfig("config.json")
	if err != nil {
		fmt.Println("read config error\n", err)
		return
	}
	fp := gofeed.NewParser()
	for {
		oldArticles, err := readOldArticles()
		if err != nil {
			fmt.Println("read old articles error\n", err)
			return
		}
		for _, source := range sources.Sources {
			feed, _ := fp.ParseURL(source.Url)
			for _, item := range feed.Items {
				if !isNew(item.Link, oldArticles) {
					continue
				}
				deprecateArticle(item.Link)
				if aiFilter(*llm, itemToArticle(*item), source) {
					fmt.Println("Article passed!\n", item.Link)
				}
			}
		}
		{
			fmt.Println("Alle articles filtered, trying again in a bit...")
			time.Sleep(time.Second * 10)
		}
	}
}

func itemToArticle(item gofeed.Item) Article {
	article := Article{
		Title:       item.Title,
		Description: item.Description,
		Url:         item.Link,
		Content:     item.Content,
	}
	return article
}

func readOldArticles() (OldArticles, error) {

	if !FileExists(oldArticlesFileName) {
		file, _ := os.Create(oldArticlesFileName)
		jsonBytes, _ := json.Marshal(OldArticles{Urls: []string{}})
		file.WriteString(string(jsonBytes))
	}
	jsonFile, err := os.Open(oldArticlesFileName)
	if err != nil {
		return OldArticles{}, err
	}
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return OldArticles{}, err
	}
	var oldArticles OldArticles
	err = json.Unmarshal(byteValue, &oldArticles)
	if err != nil {
		return OldArticles{}, err
	}
	return oldArticles, nil
}
