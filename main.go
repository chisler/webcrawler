package main

import (
	"./crawl"
	"./fetch"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
	"syscall"
)

func cmd_usage() {
	fmt.Fprintf(os.Stderr, "Usage: main http://example.com/\n")
	os.Exit(2)
}

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

func main() {
	//Receive command line arguments
	flag.Parse()
	args := flag.Args()
	fmt.Println(args)

	if len(args) < 1 {
		cmd_usage()
		fmt.Println("Please specify start page")
		os.Exit(1)
	}
	startUrl := args[0]

	//Initialize a crawler
	crawler := crawl.NewWebCrawler()

	//Initialize a fetcher
	fetcher, err := fetch.NewPageFetcher(startUrl)
	if err != nil {
		log.Fatal("Failed to parse starting URL: %s", startUrl)
		return
	}

	//Crawl pages and
	start := time.Now()
	crawler.Crawl(startUrl, 1, fetcher)
	end := time.Since(start)

	//Show results
	printResultMap(crawler, start, end)
}
