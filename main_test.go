package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestNewPageFetcher(t *testing.T) {
	_, err := NewPageFetcher("http://example.com/")
	assert.Equal(t, nil, err, "No error should be raised")

	_, err = NewPageFetcher("##$$%%%/")

	assert.NotEqual(t, nil, err, "Error should be raised")

}

func TestNormalizeUrl(t *testing.T) {
	var NormalizeUrlTests = []struct {
		input    string // input
		expected string // expected result
	}{
		{"/#fragment", "http://example.com/"},
		{"/?q=abc", "http://example.com/?q=abc"},
		{"http://otherdomain.com", "http://otherdomain.com"},
	}

	fetcher, _ := NewPageFetcher("http://example.com/")

	for _, tt := range NormalizeUrlTests {
		actual := fetcher.normalizeUrl(tt.input).String()
		assert.Equal(t, tt.expected, actual)
	}

}

func TestGetUrlsFromTags(t *testing.T) {
	startUrlString := "http://example.com/"
	fetcher, _ := NewPageFetcher(startUrlString)
	var html = `
		<a href="http://example.com/about"></a>		 	 OK
		<a href="/relative"></a>			 	 OK
		<a href="#thiswillbecomehost"></a>		  	 OK
		<a href="http://facebook.com/"></a>		 	 OK
		<a href="&*!@)$&@%(will not parse it)"></a>	 	 NOT OK
		`
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))

	allLinks := fetcher.getUrlsFromTags("a", "href", doc)

	assert.Equal(t, "http://example.com/about", allLinks[0].String())
	assert.Equal(t, "http://example.com/relative", allLinks[1].String())
	assert.Equal(t, startUrlString, allLinks[2].String())
	assert.Equal(t, "http://facebook.com/", allLinks[3].String())
	assert.Equal(t, 4, len(allLinks))
}

func TestFilterLinksByDomain(t *testing.T) {
	startUrlString := "http://example.com/"
	fetcher, _ := NewPageFetcher(startUrlString)
	var html = `
		<a href="http://example.com/about"></a>			 OK
		<a href="http://facebook.com/"></a>			 NOT OK
		`
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
	allLinks := fetcher.getAllLinks(doc)
	filteredLinks := fetcher.excludeExternalLinks(allLinks)

	assert.Equal(t, "http://example.com/about", filteredLinks[0].String())
	assert.Equal(t, 1, len(filteredLinks))
}

func TestGetStaticAssets(t *testing.T) {
	startUrlString := "http://example.com/"
	fetcher, _ := NewPageFetcher(startUrlString)
	var html = `
		<script src="https://example.com"></script> OK
		<script></script>					 NOT OK
		<script src="!@#$%"></script>				 NOT OK

		<img src="https://example.com/image">			 OK
		<img>							 NOT OK
		<img src="!@#$%">					 NOT OK

		<link rel="stylesheet" href="https://example.com/css">	 OK
		<link rel="license" href="https://example.com/license">	 OK
		<link rel="icon" href="https://example.com/icon">	 OK
		<link rel="shortcut icon" 				 OK
		href="https://example.com/sicon">

		<link rel="profile" href="http://example.com">		 NOT OK
		<link rel="dns-prefetch" href="//fonts.googleapis.com">  NOT OK

		<a href="http://example.com/about"></a>
		<a href="http://facebook.com/"></a>
		`
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
	assets := fetcher.getStaticAssets(doc)

	assert.Equal(t, "https://example.com", assets[0].String())
	assert.Equal(t, "https://example.com/image", assets[1].String())
	assert.Equal(t, "https://example.com/css", assets[2].String())
	assert.Equal(t, "https://example.com/license", assets[3].String())
	assert.Equal(t, "https://example.com/icon", assets[4].String())
	assert.Equal(t, "https://example.com/sicon", assets[5].String())
	assert.Equal(t, 6, len(assets))
}
