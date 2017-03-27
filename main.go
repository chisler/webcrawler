package main


import (
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
)

func printLinks(urlString string) {
	doc, err := goquery.NewDocument(urlString)
	if err != nil {
		log.Fatal(err)

	}

	doc.Find("a").Each(func(index int, linkTag *goquery.Selection) {
		link, _ := linkTag.Attr("href")
		fmt.Printf("Link #%d:'%s'\n", index, link)
	})
	return
}

func main() {
	printLinks("http://example.com")
}