package rss_parser

import (
	"fmt"
	"sync"
	"time"
)

type RssItem struct{
	Title string
	Source string
	SourceURL string
	Link string
	PublishDate time.Time
	Description string
}

type rssParserSync struct {
	wg *sync.WaitGroup
	out chan RssItem
}

func (p *rssParserSync) parse(url string) RssItem{
	fmt.Printf("Got url=%s\n", url)

	return RssItem{
		Title:       "",
		Source:      "",
		SourceURL:   "",
		Link:        "",
		PublishDate: time.Time{},
		Description: "",
	}
}

func newRssParserSynchronizer(out chan RssItem) *rssParserSync {
	var wg sync.WaitGroup
	return &rssParserSync{
		wg: &wg,
		out: out,
	}
}

func (p *rssParserSync) parseAsync(urls []string) {
	for _, url := range urls {
		p.wg.Add(1)
		url := url
		go func() {
			defer p.wg.Done()
			p.parse(url)
		}()
	}
	p.wg.Wait()
}

func Parse(urls []string) []RssItem {
	out := make(chan RssItem, len(urls))
	ps := newRssParserSynchronizer(out)
	ps.parseAsync(urls)
	return []RssItem{}
}
