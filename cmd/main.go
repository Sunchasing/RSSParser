package main

import (
	rss_parser "rss-reader/pkg/rss_parser_async"
)

func main(){
	rss_parser.Parse([]string{"url1", "url2"})
}
