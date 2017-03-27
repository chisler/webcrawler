package main

import (
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"
	"./fetch"
)

type Crawler interface {
	// Crawl uses fetcher to recursively and asynchronously crawl
	// pages starting with url, to a maximum of depth.
	Crawl(currentUrl string, depth int, fetcher fetch.Fetcher)
}

//WebCrawler instance - abuser of Crawler interface
type WebCrawler struct {
	Fetched map[string][]*url.URL
	sync.RWMutex
	errors  []error
}

//Constructor of WebCrawler struct
func NewWebCrawler() *WebCrawler {

	return &WebCrawler{
		Fetched: make(map[string][]*url.URL),
	}

}

//Recursively crawls pages
// extracts static assets and URLs
func (crawler *WebCrawler) Crawl(currentUrl string, depth int, fetcher fetch.Fetcher) {

	//Do not Crawl page if depth was exceeded
	if depth == 0 {
		return
	}

	//Do not Crawl page twice
	crawler.RLock()
	if _, ok := crawler.Fetched[currentUrl]; ok {
		crawler.RUnlock()
		return
	}
	crawler.RUnlock()

	//Fetch page's static assets and URLs
	assets, urls, err := fetcher.Fetch(currentUrl)
	if err != nil {
		log.Println(err)
		crawler.Lock()
		crawler.errors = append(crawler.errors, err)
		crawler.Unlock()
		return
	}

	//Add fetched resources to a storage synchronously
	crawler.Lock()
	crawler.Fetched[currentUrl] = assets
	crawler.Unlock()

	//Recursively and asynchronously Crawl URLs just fetched
	done := make(chan bool)
	for _, u := range urls {
		go func(currentUrl *url.URL) {
			crawler.Crawl(currentUrl.String(), depth-1, fetcher)
			done <- true
		}(u)
	}

	//Wait for deeper Crawls to end
	for i := 0; i < len(urls); i++ {
		<-done
	}

	return
}

func main() {
	startUrl := "http://jeniasofronov.ru/"
	crawler := NewWebCrawler()
	fetcher, err := fetch.NewPageFetcher(startUrl)


	if err != nil {
		log.Fatal("Failed to parse starting URL: %s", startUrl)
		return
	}


	start := time.Now()

	crawler.Crawl(startUrl, 3, fetcher)
	fmt.Println("_________MAP__________")
	for k, v := range crawler.Fetched {
		fmt.Printf("Node: %v \nAssets: %v \n\n", k, v)
	}
	fmt.Printf("Extracted %v pages\n", len(crawler.Fetched))

	end := time.Since(start)
	fmt.Printf("Execution took %v", end)

}
