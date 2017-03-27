/*
Package crawl implements a simple library for crawling web site with max depth

Depends on fetch package.

Usage:

	crawler := crawl.NewWebCrawler()
	fetcher, err := fetch.NewPageFetcher(startUrl)

	crawler.Crawl(startUrl, depth, fetcher)

*/

package crawl

//Inspired by https://tour.golang.org/concurrency/10

import (
	"../fetch"
	"log"
	"net/url"
	"sync"
)

type Crawler interface {
	// Crawl uses fetcher to recursively and asynchronously crawl
	// pages starting with url, to a maximum of depth.
	Crawl(currentUrl string, depth int, fetcher fetch.Fetcher)
}

//WebCrawler instance - abuser of Crawler interface
type WebCrawler struct {
	Fetched map[string]*Page
	errors  []error

	sync.RWMutex
}

//Constructor of WebCrawler struct
func NewWebCrawler() *WebCrawler {

	return &WebCrawler{
		Fetched: make(map[string]*Page),
	}

}

type Page struct {
	Assets []*url.URL
	Urls   []*url.URL
}

func NewPage(assets []*url.URL, urls []*url.URL) *Page {
	return &Page{
		Assets: assets,
		Urls:   urls,
	}
}

//Recursively crawls pages
// extracts static Assets and URLs
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

	//Fetch page's static Assets and URLs
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
	crawler.Fetched[currentUrl] = NewPage(assets, urls)
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
