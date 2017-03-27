package main


import (
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
)



func printStaticAssets(urlString string) {
	doc, err := goquery.NewDocument(urlString)
	if err != nil {
		log.Fatal(err)

	}

	doc.Find("script").Each(func(index int, linkTag *goquery.Selection) {
		link, _ := linkTag.Attr("src")
		fmt.Printf("Script #%d:'%s'\n", index, link)
	})

	doc.Find("img").Each(func(index int, linkTag *goquery.Selection) {
		link, _ := linkTag.Attr("src")
		fmt.Printf("Img #%d:'%s'\n", index, link)
	})

	doc.Find("link").Each(func(index int, linkTag *goquery.Selection) {
		link, _ := linkTag.Attr("href")
		fmt.Printf("Link #%d:'%s'\n", index, link)
	})
}

func printLinks(urlString string) {
	doc, err := goquery.NewDocument(urlString)
	if err != nil {
		log.Fatal(err)

	}

	doc.Find("a").Each(func(index int, linkTag *goquery.Selection) {
		link, _ := linkTag.Attr("href")
		fmt.Printf("A #%d:'%s'\n", index, link)
	})
}

func main() {
	printLinks("https://monzo.com/")
	printStaticAssets("https://monzo.com/")
}