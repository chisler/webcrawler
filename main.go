package main

import (
	"fmt"
	"./fetch"
)

func main() {
	fetcher, _ := fetch.NewPageFetcher("https://monzo.com/")
	fmt.Print(fetcher.Fetch("https://monzo.com/"))
}

