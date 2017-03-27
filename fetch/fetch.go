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
	Fetch(targetUrl string) (assets []*url.URL, urls []*url.URL, err error)
}

type PageFetcher struct {
	startUrl *url.URL
}

func NewPageFetcher(startUrlString string) (*PageFetcher, error) {
	startUrl, err := url.Parse(startUrlString)

	if err != nil {
		return nil, err
	}

	return &PageFetcher{startUrl: startUrl}, nil

}

func (f *PageFetcher) Fetch(targetUrl string) (assets []*url.URL, urls []*url.URL, err error) {

	doc, err := goquery.NewDocument(targetUrl)
	if err != nil {
		return nil, nil, err
	}

	urls = f.getInternalLinks(doc)

	assets = f.getStaticAssets(doc)

	return
}

func (f *PageFetcher) getInternalLinks(doc *goquery.Document) (res []*url.URL) {

	allLinks := f.getAllLinks(doc)
	res = f.excludeExternalLinks(allLinks)

	return
}

func (f *PageFetcher) getAllLinks(doc *goquery.Document) (res []*url.URL) {

	res = f.getUrlsFromTags("a", "href", doc)
	return
}

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

func (f *PageFetcher) excludeExternalLinks(urls []*url.URL) (filteredLinks []*url.URL) {

	filteredLinks = urls[:0]
	for _, currentUrl := range urls {
		if f.startUrl.Host == currentUrl.Host {
			filteredLinks = append(filteredLinks, currentUrl)
		}
	}

	return
}

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
