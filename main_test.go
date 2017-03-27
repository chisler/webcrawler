package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"net/url"
	"github.com/PuerkitoBio/goquery"
	"strings"
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
