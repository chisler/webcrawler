package main


import (
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
)


func Fetch(urlString string)  {
	doc, err := goquery.NewDocument(urlString)
	if err != nil {
		log.Fatal(err)

	}

	fmt.Println(getLinks(doc))
	fmt.Println(getStaticAssets(doc))
}

func getStaticAssets(doc *goquery.Document) (res []string) {

	//Add <script> tag assets
	res = getAttrsFromTags(doc, "script", "src")

	//Add <img> tag assets
	res = append(res, getAttrsFromTags(doc, "img", "src")...)

	//Add <link> tag assets
	res = append(res, getAttrsFromTags(doc, "link", "href")...)
	return
}

func getLinks(doc *goquery.Document) (res []string) {
	res = getAttrsFromTags(doc, "a", "href")
	return
}

//Returns attrs from document by tag and attr
func getAttrsFromTags(doc *goquery.Document, tagName, attrName string) (res []string) {

	doc.Find(tagName).Each(func(index int, linkTag *goquery.Selection) {

		if link, ok := linkTag.Attr(attrName); ok {
			res = append(res, link)
		}
	})
	return
}

func main() {
	Fetch("https://monzo.com/")
}