package main

import (
	"./crawl"
	"./fetch"
	"fmt"
	"log"
	"time"
)

func main() {
	startUrl := "http://monzo.com/"

	crawler := crawl.NewWebCrawler()
	fetcher, err := fetch.NewPageFetcher(startUrl)

	if err != nil {
		log.Fatal("Failed to parse starting URL: %s", startUrl)
		return
	}

	start := time.Now()

	crawler.Crawl(startUrl, 2, fetcher)
	fmt.Println("_________MAP__________")
	for k, v := range crawler.Fetched {
		fmt.Printf("Node: %v \nAssets: %v \n\n", k, v)
	}
	fmt.Printf("Extracted %v pages\n", len(crawler.Fetched))

	end := time.Since(start)
	fmt.Printf("Execution took %v", end)
}
