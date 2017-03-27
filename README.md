# webcrawler
Library for crawling web site with depth limit.

## Getting started

### How it works

Given a root URL, it returns a site map, showing which static assets each page depends on, and the links between pages.

The crawler is limited to one domain - crawling example.com it crawls all pages within the domain,
but does not follow external links.

### Installing

Installations requires Go1.1 for goquery.
```
$ go get github.com/chisler/webcrawler
```

### Usage


Depth has its default, startUrl does not.
```

$ webcrawler 
Usage: main -startUrl=http://example.com/ -depth=3
  -depth int
        Depth of crawling. (default 4)
  -startUrl string
        Root URL of website to crawl.
```

Proper running 
```
$ webcrawler -startUrl=http://monzo.com -depth=2
...
$ ls
result.txt
```
$GOPATH should be [set](https://golang.org/doc/code.html#GOPATH).

What's in the box
```
$ head result.txt 
_________MAP__________

Extracted 14 pages
Execution started at 2017-03-27 15:50:22.692179718 +0300 MSK
Execution took 3.431152137s

Node: http://monzo.com/-play-store-redirect 
Urls: [urllist] 
Assets: [assetlist] 
```

Run tests
```
$ go test github.com/chisler/webcrawler/crawl
ok      github.com/chisler/webcrawler/crawl     0.003s
$ go test github.com/chisler/webcrawler/fetch
ok      github.com/chisler/webcrawler/fetch     0.008s
```

## Built With

* [golang](https://golang.org/doc/install)
* [goquery](https://github.com/PuerkitoBio/goquery)
* [testify](https://github.com/stretchr/testify/assert)
