package crawl

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

const testStartUrl = "http://golang.org/"

func TestWebCrawler_Crawl(t *testing.T) {
	crawler := NewWebCrawler()

	crawler.Crawl(testStartUrl, 3, fetcher)

	assert.Equal(t, 4, len(crawler.Fetched))

	//Some pages could not be found
	assert.NotNil(t, crawler.errors)

	//Assert that crawler has processed everything available
	for k := range crawler.Fetched {
		assert.NotNil(t, fetcher[k])
	}

}

func TestWebCrawler_Page(t *testing.T) {
	crawler := NewWebCrawler()

	crawler.Crawl(testStartUrl, 1, fetcher)
	expected_urls := parseUrls(fetcher[testStartUrl].urls)

	assert.Equal(t, expected_urls, crawler.Fetched[testStartUrl].Urls)
	assert.Nil(t, crawler.Fetched[testStartUrl].Assets)
}

func TestWebCrawler_CrawlDepth(t *testing.T) {
	crawler := NewWebCrawler()

	crawler.Crawl(testStartUrl, 1, fetcher)

	assert.Equal(t, 1, len(crawler.Fetched))
	assert.Nil(t, crawler.errors)
}

func TestWebCrawler_FetchErrorCallback(t *testing.T) {
	crawler := NewWebCrawler()

	crawler.Crawl("No such url", 1, fetcher)

	assert.Equal(t, 0, len(crawler.Fetched))
	assert.NotNil(t, crawler.errors)
}

//------------Mock for Fetcher interface--------------
//Inspired by https://tour.golang.org/concurrency/10

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	assets []string
	urls   []string
	err    error
}

// Imitates the Fetcher interface Fetch (Assets, Urls, err)
func (f fakeFetcher) Fetch(targetUrl string) (assets []*url.URL, urls []*url.URL, err error) {
	if res, ok := f[targetUrl]; ok {
		return parseUrls(res.assets), parseUrls(res.urls), nil
	}
	return nil, nil, fmt.Errorf("not found: %s", targetUrl)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"http://golang.org/": &fakeResult{
		[]string{},
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
		nil,
	},
	"http://golang.org/pkg/": &fakeResult{
		[]string{},
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
		nil,
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		[]string{},
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
		nil,
	},
	"http://golang.org/pkg/os/": &fakeResult{
		[]string{},
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
		nil,
	},
}

//---------------------Helpers------------------------

func parseUrls(urlsString []string) (urls []*url.URL) {
	for _, urlString := range urlsString {
		if parsed, err := url.Parse(urlString); err == nil {
			urls = append(urls, parsed)
		}
	}
	return
}
