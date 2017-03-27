package main

import (
	"log"
	"time"
	"fmt"
	"./crawl"
	"./fetch"
	"os"
	"flag"
)

func cmd_usage() {
	fmt.Fprintf(os.Stderr, "Usage: main http://example.com/\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Parse()

	args := flag.Args()
	fmt.Println(args)
	if len(args) < 1 {
		cmd_usage()
		fmt.Println("Please specify start page")
		os.Exit(1)
	}

	startUrl := args[0]

	crawler := crawl.NewWebCrawler()
	fetcher, err := fetch.NewPageFetcher(startUrl)


	if err != nil {
		log.Fatal("Failed to parse starting URL: %s", startUrl)
		return
	}

	start := time.Now()
	crawler.Crawl(startUrl, 3, fetcher)
	end := time.Since(start)

	fmt.Println("_________MAP__________")
	for k, v := range crawler.Fetched {
		fmt.Printf("Node: %v \nUrls: %v \nAssets: %v \n\n", k, v.Urls, v.Assets)
	}
	fmt.Printf("Extracted %v pages\n", len(crawler.Fetched))

	fmt.Printf("Execution took %v", end)
}
