package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mmcdole/gofeed"
	"github.com/valyala/fasthttp"
)

type MangaEntry struct {
	OldId     string
	NewId     int64
	Title     string
	Chapter   string
	ScanGroup string
	Link      string
}

func rssFetchFeed(url string) (*gofeed.Feed, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(url)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := client.Do(req, resp); err != nil {
		return nil, err
	}

	body := resp.Body()
	fp := gofeed.NewParser()
	feed, err := fp.ParseString(string(body))
	if err != nil {
		return nil, err
	}
	return feed, nil
}

func rssGetLatest() (*gofeed.Feed, error) {
	url := "https://api.mangaupdates.com/v1/releases/rss"
	feed, err := rssFetchFeed(url)
	if err != nil {
		feed, err = rssFetchFeed(url)
		if err != nil {
			return nil, fmt.Errorf("Could not get MangaUpdates RSS feed")
		}
	}
	return feed, nil
}

func RssParseFeed() ([]MangaEntry, error) {
	feed, err := rssGetLatest()
	if err != nil {
		return nil, err
	}

	mangaList := []MangaEntry{}
	chapterRegex := regexp.MustCompile(`(v\.\d+\s)?c\.\d+[a-z]?(\.\d+)?(-\d+(\.\d+)?)?`)
	scanGroupRegex := regexp.MustCompile(`\[(.*?)\]`)

	for _, entry := range feed.Items {
		title := entry.Title
		chapter := ""
		scanGroup := ""
		link := entry.Link

		chapterMatch := chapterRegex.FindString(title)
		if chapterMatch != "" {
			title = strings.TrimSpace(strings.Replace(title, chapterMatch, "", -1))
			chapter = chapterMatch
		}

		scanGroupMatch := scanGroupRegex.FindStringSubmatch(title)
		if len(scanGroupMatch) > 1 {
			title = strings.TrimSpace(strings.Replace(title, scanGroupMatch[0], "", -1))
			scanGroup = scanGroupMatch[1]
		}

		if strings.HasPrefix(link, "http://") {
			link = "https://" + strings.TrimPrefix(link, "http://")
		}

		mangaList = append(mangaList, MangaEntry{
			Title:     title,
			Chapter:   chapter,
			ScanGroup: scanGroup,
			Link:      link,
		})
	}
	return mangaList, nil
}
