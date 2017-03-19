package main

import (
	"image"
	"log"
	"net/http"
	"time"

	"github.com/nlopes/slack"
)

func DownloadImage(file *slack.File, slackToken string) image.Image {
	client := &http.Client{
		Timeout: time.Second * 20,
	}
	request, err := http.NewRequest(http.MethodGet, file.URLPrivateDownload, nil)
	if err != nil {
		log.Fatalf("error creating request: %s", err)
	}
	request.Header.Add("Authorization", "Bearer "+slackToken)
	response, err := client.Do(request)
	if err != nil {
		log.Fatalf("error downloading file\n%v\n%v", file, err)
	}
	defer response.Body.Close()
	downloadImg, _, err := image.Decode(response.Body)
	if err != nil {
		log.Fatalf("error downloading file\n%v\n%v", file, err)
	}
	return downloadImg
}
