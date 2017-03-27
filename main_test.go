package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"net/url"
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