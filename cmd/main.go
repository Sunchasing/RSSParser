package main

import (
	"fmt"
	rss_parser "rss-reader/pkg/rss_parser_async"
)

func main() {
	urls := []string{"https://time.com/feed/", "https://medium.com/@vaidehijoshi/feed", "https://medium.com/feed/@vaidehijoshi"}
	parsed, err := rss_parser.Parse(urls)
	if err != nil {
		return
	}
	fmt.Println(parsed)
}
