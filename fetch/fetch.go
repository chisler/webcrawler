/*
Package fetch implements a simple library for fetching web page static assets and links by URL.

Usage:

	fetcher, err := fetch.NewPageFetcher("https://startUrl.com/")

	assets, urls, err := fetcher.Fetch("https://fetchUrl.com/")

*/

package fetch

import (
	"net/url"
	"github.com/PuerkitoBio/goquery"
)

//Types "n" where <link type="n"> relates to a static asset
var linkTagAllowedTypes = map[string]bool{
	"license":       true,
	"stylesheet":    true,
	"icon":          true,
	"shortcut icon": true,
}

type Fetcher interface {
	// Fetch returns two slices: static assets and
	// URLs found on the page by its URL and
	Fetch(targetUrl string) (assets []*url.URL, urls []*url.URL, err error)
}

//PageFetcher instance - abuser of Fetcher interface
type PageFetcher struct {
	startUrl *url.URL
}

//Constructor of PageFetcher struct
func NewPageFetcher(startUrlString string) (*PageFetcher, error) {
	startUrl, err := url.Parse(startUrlString)

	if err != nil {
		return nil, err
	}

	return &PageFetcher{startUrl: startUrl}, nil

}

// Returns static assets and URLs from retrieved webpage
func (f *PageFetcher) Fetch(targetUrl string) (assets []*url.URL, urls []*url.URL, err error) {

	doc, err := goquery.NewDocument(targetUrl)
	if err != nil {
		return nil, nil, err
	}

	//Both urls and assets are allowed to be empty
	urls = f.getInternalLinks(doc)

	assets = f.getStaticAssets(doc)

	return
}

//Returns all links limited to one domain of startUrl
func (f *PageFetcher) getInternalLinks(doc *goquery.Document) (res []*url.URL) {

	allLinks := f.getAllLinks(doc)
	res = f.excludeExternalLinks(allLinks)

	return
}

//Returns all links of the page
func (f *PageFetcher) getAllLinks(doc *goquery.Document) (res []*url.URL) {

	res = f.getUrlsFromTags("a", "href", doc)
	return
}

//Returns all static assets of the page
func (f *PageFetcher) getStaticAssets(doc *goquery.Document) (res []*url.URL) {

	//Add <script> tag assets
	scriptSources := f.getUrlsFromTags("script", "src", doc)
	res = append(res, scriptSources...)

	//Add <img> tag assets
	imgSources := f.getUrlsFromTags("img", "src", doc)
	res = append(res, imgSources...)

	//Add <link> tag assets
	doc.Find("link").Each(func(_ int, linkTag *goquery.Selection) {
		if urlAttr, ok := linkTag.Attr("href"); ok && urlAttr != "" {
			absoluteUrl := f.normalizeUrl(urlAttr)
			if absoluteUrl != nil {

				//Check if link relates to a static asset
				if t, ok := linkTag.Attr("rel"); ok && linkTagAllowedTypes[t] {
					res = append(res, absoluteUrl)
				}
			}
		}
	})

	return
}

//Returns URLs from document
//by tag and attr of the url
func (f *PageFetcher) getUrlsFromTags(tagName, attrName string, doc *goquery.Document) (res []*url.URL) {

	doc.Find(tagName).Each(func(_ int, linkTag *goquery.Selection) {
		if urlAttr, ok := linkTag.Attr(attrName); ok && urlAttr != "" {

			absoluteUrl := f.normalizeUrl(urlAttr)
			if absoluteUrl != nil {
				res = append(res, absoluteUrl)
			}
		}
	})
	return
}

//Returns URLs filtered by host of PageFetcher.startUrl
func (f *PageFetcher) excludeExternalLinks(urls []*url.URL) (filteredLinks []*url.URL) {

	filteredLinks = urls[:0]
	for _, currentUrl := range urls {
		if f.startUrl.Host == currentUrl.Host {
			filteredLinks = append(filteredLinks, currentUrl)
		}
	}

	return
}

//Returns absolute URL without fragment
func (f *PageFetcher) normalizeUrl(urlString string) (normalizedUrl *url.URL) {

	// Parse and resolve to an absolute url
	normalizedUrl, err := f.startUrl.Parse(urlString)
	if err != nil {
		return nil
	}

	// Remove fragment
	normalizedUrl.Fragment = ""

	return
}
