package main

import (
	"fmt"
	"log"
	rss_parser "rss-reader/pkg/rss_parser_async"
)

func main() {
	urls := []string{"https://time.com/feed/", "https://medium.com/feed/@vaidehijoshi"}
	parsed, err := rss_parser.Parse(urls)
	if err != nil {
		log.Fatalf("encountered error in RSS Parser: %v", err)
	}

	for _, r := range parsed {
		fmt.Println("----- ITEM -----")
		fmt.Printf("Title=%v\n", r.Title)
		fmt.Printf("Source=%v\n", r.Source)
		fmt.Printf("SourceURL=%v\n", r.SourceURL)
		fmt.Printf("Link=%v\n", r.Link)
		fmt.Printf("PublishDate=%v\n", r.PublishDate)
		fmt.Printf("Description=%v\n\n", r.Description)
	}
}
