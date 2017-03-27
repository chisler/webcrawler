package main


import (
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
	"net/url"
)

//Types "n" where <link type="n"> relates to a static asset
var linkTagAllowedTypes = map[string]bool{
	"license":       true,
	"stylesheet":    true,
	"icon":          true,
	"shortcut icon": true,
}


func Fetch(urlString string)  {
	doc, err := goquery.NewDocument(urlString)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(getLinks(doc))
	fmt.Println(getStaticAssets(doc))
}

func getStaticAssets(doc *goquery.Document) (res []*url.URL) {

	//Add <script> tag assets
	res = getAttrsFromTags(doc, "script", "src")

	//Add <img> tag assets
	res = append(res, getAttrsFromTags(doc, "img", "src")...)

	//Add <link> tag assets
	doc.Find("link").Each(func(_ int, linkTag *goquery.Selection) {
		if urlAttr, ok := linkTag.Attr("href"); ok && urlAttr != "" {
			if absoluteUrl := normalizeUrl(urlAttr, doc.Url); absoluteUrl != nil {

				//Check if link relates to a static asset
				if t, ok := linkTag.Attr("rel"); ok && linkTagAllowedTypes[t] {
					res = append(res, absoluteUrl)
				}
			}
		}
	})
	return
}

func getLinks(doc *goquery.Document) (res []*url.URL) {
	res = getAttrsFromTags(doc, "a", "href")
	return
}

//Returns attrs from document by tag and attr
func getAttrsFromTags(doc *goquery.Document, tagName, attrName string) (res []*url.URL) {

	doc.Find(tagName).Each(func(index int, linkTag *goquery.Selection) {
		if urlAttr, ok := linkTag.Attr(attrName); ok && urlAttr != "" {
			if absolute := normalizeUrl(urlAttr, doc.Url); absolute != nil {
				res = append(res, absolute)
			}
		}
	})
	return
}

func normalizeUrl(urlString string, host *url.URL) (normalizedUrl *url.URL) {

	normalizedUrl, err := host.Parse(urlString)
	if err != nil {
		return nil
	}

	return
}


func main() {
	Fetch("https://monzo.com/")
}
