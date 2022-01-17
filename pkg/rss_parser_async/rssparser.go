package rss_reader

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

// region XMLParser
type rss struct {
	XMLChannel xml.Name `xml:"rss"`
	Channel    channel  `xml:"channel"`
}

type channel struct {
	XMLName        xml.Name      `xml:"channel"`
	XMLDescription []description `xml:"description"`
	XMLItem        []item        `xml:"item"`
	Title          string        `xml:"title"`
}

type item struct {
	Title       string `xml:"title"`
	Source      string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
}

type description struct {
	XMLName xml.Name `xml:"description"`
	XMLLink string   `xml:"link"`
}

func parseTime(timeStr string) (*time.Time, error) {

	timeFormats := []string{
		"01/02 03:04:05PM '06 -0700",          //Layout
		"Mon Jan _2 15:04:05 2006",            //ANSI
		"Mon Jan _2 15:04:05 MST 2006",        //UnixDate
		"Mon Jan 02 15:04:05 -0700 2006",      //RubyDate
		"02 Jan 06 15:04 MST",                 //RFC822
		"02 Jan 06 15:04 -0700",               //RFC822Z
		"Monday, 02-Jan-06 15:04:05 MST",      //RFC850
		"Mon, 02 Jan 2006 15:04:05 MST",       //RFC1123
		"Mon, 02 Jan 2006 15:04:05 -0700",     //RFC1123Z
		"2006-01-02T15:04:05Z07:00",           //RFC3339
		"2006-01-02T15:04:05.999999999Z07:00", //RFC3339Nano
	}
	for _, tf := range timeFormats {

		parsedTime, err := time.Parse(tf, timeStr)
		if err == nil {
			return &parsedTime, nil
		}
	}
	return nil, fmt.Errorf("could not parse %v into any time.Time format", timeStr)
}

const regex = `<.*?>`

func stripHtmlRegex(s string) string {
	//todo: avoid compiling at each check + error handling
	r := regexp.MustCompile(regex)
	return r.ReplaceAllString(s, "")
}

func parseXml(xmlBytes []byte, url string) (rv []RssItem, err error) {
	var rssHolder rss
	err = xml.Unmarshal(xmlBytes, &rssHolder)
	if err != nil {
		err = fmt.Errorf("could not unmarshal xml. Error was: %v", err)
		return
	}

	for _, elem := range rssHolder.Channel.XMLItem {
		parsedTime, e := parseTime(elem.PubDate)
		if e != nil {
			err = fmt.Errorf("error while parsing publishing date: %v", e)
			return
		}
		newItem := RssItem{
			Title:       strings.TrimSpace(elem.Title),
			Source:      rssHolder.Channel.Title,
			SourceURL:   url,
			Link:        strings.TrimSpace(elem.Source),
			PublishDate: *parsedTime,
			Description: stripHtmlRegex(strings.TrimSpace(elem.Description)),
		}
		rv = append(rv, newItem)
	}
	return rv, nil
}

// endregion

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

func (p *rssParserSync) parseToChannel(url string) (err error) {
	defer p.wg.Done()
	var rv []RssItem
	bytesArr, err := p.hr.get(url)
	if err != nil {
		err = fmt.Errorf("failed to GET feed from %s: %v", url, err)
		return
	} else {

		parsedXml, e := parseXml(bytesArr, url)
		if e != nil {
			err = fmt.Errorf("failed to parse XML data: %v", e)
			return
		}
		for _, parsedItem := range parsedXml {
			rv = append(rv, parsedItem)

		}

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
			if err != nil {
				err = fmt.Errorf("failed to get return value for %s: %v", urlLocal, err)
			}
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
