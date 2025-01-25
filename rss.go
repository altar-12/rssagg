package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string     `xml:"title"`
		Link        string     `xml:"link"`
		Description string     `xml:"description"`
		Language    string     `xml:"language"`
		Items       []FeedItem `xml:"item"`
	} `xml:"channel"`
}

type FeedItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
}

func urlToRSSFeed(url string) (RSSFeed, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return RSSFeed{}, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return RSSFeed{}, err
	}
	feed := RSSFeed{}
	err = xml.Unmarshal(data, &feed)
	if err != nil {
		return RSSFeed{}, err
	}
	return feed, nil
}
