package main

import (
	"github.com/chisler/webcrawler/crawl"
	"github.com/chisler/webcrawler/fetch"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
	"syscall"
)

var (
	startUrl = flag.String("startUrl", "", "Root URL of website to crawl.")
	depth    = flag.Int("depth", 4, "Depth of crawling.")
)

func main() {
	//Receive command line arguments
	flag.Parse()

	//Check cmd arguments
	if *startUrl == "" {
		cmd_usage()
	}

	//Initialize a crawler
	crawler := crawl.NewWebCrawler()

	//Initialize a fetcher
	fetcher, err := fetch.NewPageFetcher(*startUrl)
	if err != nil {
		log.Fatal("Failed to parse starting URL: %s", startUrl)
		return
	}

	//Crawl pages
	start := time.Now()
	crawler.Crawl(*startUrl, *depth, fetcher)
	end := time.Since(start)

	//Show results
	printResultMap(crawler, start, end)
}

// Prints the sitemap to textfile or stdout in plain text
func printResultMap(crawler *crawl.WebCrawler, start time.Time, end time.Duration) {
	f, err := os.Create("result.txt")

	if err != nil {
		log.Panic(err)
		log.Println("Result will be printed to stdout.")
		f = os.NewFile(uintptr(syscall.Stdout), "/dev/stdout")
	}

	f.WriteString("_________MAP__________\n\n")

	f.WriteString(fmt.Sprintf("Extracted %v pages\n", len(crawler.Fetched)))
	f.WriteString(fmt.Sprintf("Extracted %v pages\n", len(crawler.Fetched)))
	f.WriteString(fmt.Sprintf("Execution started at %v\n", start))
	f.WriteString(fmt.Sprintf("Execution took %v\n\n", end))

	for k, v := range crawler.Fetched {
		f.WriteString(fmt.Sprintf("Node: %v \nUrls: %v \nAssets: %v \n\n", k, v.Urls, v.Assets))
	}
}

// Prints usage in case of inproper cmd arguments 
func cmd_usage() {
	fmt.Fprintf(os.Stderr, "Usage: main -startUrl=http://example.com/ -depth=3\n")
	flag.PrintDefaults()
	os.Exit(2)
}
