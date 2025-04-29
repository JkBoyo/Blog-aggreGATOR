package main

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
)

func fetchFeed(ctx context.Context, feedUrl string) (*RSSFeed, error) {
	feedReq, err := http.NewRequestWithContext(ctx, "GET", feedUrl, nil)
	if err != nil {
		return nil, err
	}

	feedReq.Header.Set("User-Agent", "gator")

	resp, err := http.DefaultClient.Do(feedReq)
	if err != nil {
		return &RSSFeed{}, err
	}
	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return &RSSFeed{}, err
	}

	var fetchedFeed RSSFeed

	err = xml.Unmarshal(dat, &fetchedFeed)
	if err != nil {
		return &RSSFeed{}, err
	}

	for i, iT := range fetchedFeed.Channel.Item {
		fetchedFeed.Channel.Item[i].Title = html.UnescapeString(iT.Title)
		fetchedFeed.Channel.Item[i].Description = html.UnescapeString(iT.Description)
	}
	fetchedFeed.Channel.Title = html.UnescapeString(fetchedFeed.Channel.Title)
	fetchedFeed.Channel.Description = html.UnescapeString(fetchedFeed.Channel.Description)

	return &fetchedFeed, nil
}

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}
