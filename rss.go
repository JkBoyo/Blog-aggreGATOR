package main

import (
	"GATOR/internal/database"
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
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

func scrapeFeeds(s *state) error {
	feedToFetch, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	fetchedFeedParams := database.MarkFeedFetchedParams{
		UpdatedAt: time.Now(),
		ID:        feedToFetch.ID,
	}
	err = s.db.MarkFeedFetched(context.Background(), fetchedFeedParams)
	if err != nil {
		return err
	}
	feed, err := fetchFeed(context.Background(), feedToFetch.Url)
	if err != nil {
		return err
	}
	fmt.Println("Feed Title: ", feed.Channel.Title)
	for _, item := range feed.Channel.Item {
		pubDateTime, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			fmt.Println(fmt.Errorf("PubDate parsing error: %v", err))
			continue
		}
		if item.Link == "" {
			fmt.Println(errors.New("no url found"))
			continue
		}
		postParams := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: item.Description != ""},
			PublishedAt: sql.NullTime{Time: pubDateTime, Valid: item.PubDate != ""},
			FeedID:      feedToFetch.ID,
		}
		post, err := s.db.CreatePost(context.Background(), postParams)
		if err != nil {
			pgErr, ok := err.(*pq.Error)
			if ok && pgErr.Code == "23505" {
				continue
			} else {
				fmt.Println(err)
				continue
			}
		}
		fmt.Println("  Post: ", post)
	}
	return nil
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
