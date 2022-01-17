package rss_reader

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type RssItem struct {
	Title       string
	Source      string
	SourceURL   string
	Link        string
	PublishDate time.Time
	Description string
}

type rssParserSync struct {
	wg  *sync.WaitGroup
	out chan []RssItem
	hr  httpInterface
}

type httpInterface interface {
	get(url string) (responseBytes []byte, err error)
}

type httpClient struct {
	httpInterface
}

func (r *httpClient) get(url string) (responseBytes []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get RSS from url %s: %v", url, err)
	}

	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("response code for %s was %v", url, resp.StatusCode)
	}

	responseBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	return
}

// Data us not being formatted yet, used for testing
func (p *rssParserSync) parseToChannel(url string) (err error) {
	defer p.wg.Done()

	_, err = p.hr.get(url)
	if err != nil {
		return fmt.Errorf("failed to GET feed from %s: %v", url, err)
	} else {

		var rv []RssItem
		rv = append(rv, RssItem{
			Title:       "",
			Source:      "",
			SourceURL:   "",
			Link:        "",
			PublishDate: time.Time{},
			Description: "",
		})

		p.out <- rv
	}
	return
}

func newRssParserSynchronizer(out chan []RssItem) *rssParserSync {
	var wg sync.WaitGroup
	return &rssParserSync{
		wg:  &wg,
		out: out,
		hr:  &httpClient{},
	}
}

func (p *rssParserSync) parseAsync(urls []string) (err error) {
	for _, url := range urls {
		p.wg.Add(1)
		urlLocal := url
		go func() {
			err = p.parseToChannel(urlLocal)
			err = fmt.Errorf("failed to get return value for %s, %v", urlLocal, err)
		}()
	}
	p.wg.Wait()
	close(p.out)
	return
}

func (p *rssParserSync) formatChannelData() (rv []RssItem) {
	for items := range p.out {
		for _, item := range items {
			rv = append(rv, item)
		}
	}
	return
}

func Parse(urls []string) (rv []RssItem, err error) {
	out := make(chan []RssItem, len(urls))
	ps := newRssParserSynchronizer(out)
	if err = ps.parseAsync(urls); err != nil {
		return nil, fmt.Errorf("failed to parseToChannel URLs: %v", err)
	}

	rv = ps.formatChannelData()
	return
}
