package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"net/url"
	"strings"
	"testing"
)

func TestNormalizeUrl(t *testing.T) {
	var NormalizeUrlTests = []struct {
		input    string // input
		expected string // expected result
	}{
		{"/?q=abc", "http://example.com/?q=abc"},
		{"http://otherdomain.com", "http://otherdomain.com"},
	}

	host, _ := url.Parse("http://example.com/")

	for _, tt := range NormalizeUrlTests {
		actual := normalizeUrl(tt.input, host).String()
		assert.Equal(t, tt.expected, actual)
	}

}

func TestGetUrlsFromTags(t *testing.T) {
	startUrlString := "http://example.com/"
	startUrl, _ := url.Parse(startUrlString)
	var html = `
		<a href="http://example.com/about"></a>		 	 OK
		<a href="/relative"></a>			 	 OK
		<a href="http://facebook.com/"></a>		 	 OK
		<a href="&*!@)$&@%(will not parse it)"></a>	 	 NOT OK
		`

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
	doc.Url = startUrl
	allLinks := getUrlsFromTags("a", "href", doc)
	fmt.Print(allLinks)

	assert.Equal(t, 3, len(allLinks))
	assert.Equal(t, "http://example.com/about", allLinks[0].String())
	assert.Equal(t, "http://example.com/relative", allLinks[1].String())
	assert.Equal(t, "http://facebook.com/", allLinks[2].String())
}

func TestExcludeExternalLinks(t *testing.T) {
	host, _ := url.Parse("http://example.com/")
	var html = `
		<a href="http://example.com/about"></a>			 OK
		<a href="http://facebook.com/"></a>			 NOT OK
		`
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
	allLinks := getAllLinks(doc)
	filteredLinks := excludeExternalLinks(allLinks, host)

	assert.Equal(t, "http://example.com/about", filteredLinks[0].String())
	assert.Equal(t, 1, len(filteredLinks))
}
