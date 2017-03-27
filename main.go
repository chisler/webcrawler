package main

import (
	"sync"
	"net/url"
	"./fetch"
	"log"
	"time"
	"fmt"
)

type Storage struct {
	fetched map[string][]*url.URL
	sync.RWMutex
}

var store Storage

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
// Based on golang tour solution
func Crawl(currentUrl string, depth int, fetcher fetch.Fetcher) {
	if depth == 0 {
		return
	}

	store.RLock()
	if _, ok := store.fetched[currentUrl]; ok {
		store.RUnlock()
		return
	}
	store.RUnlock()

	assets, urls, err := fetcher.Fetch(currentUrl)

	if err != nil {
		log.Println(err)
		return
	}

	store.Lock()
	store.fetched[currentUrl] = assets
	store.Unlock()

	done := make(chan bool)
	for _, u := range urls {
		go func(currentUrl *url.URL) {
			Crawl(currentUrl.String(), depth-1, fetcher)
			done <- true
		}(u)
	}

	for i := 0; i < len(urls); i++ {
		<-done
	}
	return
}

func main() {
	startUrl := "http://jeniasofronov.ru/"
	fetcher, err := fetch.NewPageFetcher(startUrl)

	if err != nil {
		log.Fatal("Failed to parse starting URL: %s", startUrl)
		return
	}

	store.fetched = make(map[string][]*url.URL)

	start := time.Now()

	Crawl(startUrl, 3, fetcher)

	fmt.Println("_________MAP__________")
	for k, v := range store.fetched {
		fmt.Printf("Node: %v \nAssets: %v \n\n", k, v)
	}
	end := time.Since(start)
	fmt.Printf("Extracted %v pages\n", len(store.fetched))
	fmt.Printf("Execution took %v", end)
}
