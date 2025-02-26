package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Notification struct {
	Channel string `json:"channel"`
	Title   string `json:"title"`
	Body    string `json:"body"`
	Url     string `json:"url"`
}

func sendNotification(notification Notification) error {
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}
	body, err := json.Marshal(notification)
	if err != nil {
		return err
	}
	client := &http.Client{}
	req, err := http.NewRequest(
		"POST",
		"https://api.pushify.net/v1/send",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+os.Getenv("PUSHIFY_API_KEY"))
	req.Header.Add("Content-Type", "application/json")
	fmt.Println(req)
	client.Do(req)
	return nil
}

func notificationFromArticle(article Article) Notification {
	notification := Notification{
		Title:   "SEE THIS!",
		Body:    article.Title,
		Url:     article.Url,
		Channel: "c_01jmw361knbe9fwdqafrq17nv5",
	}
	return notification
}
