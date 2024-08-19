package main

import (
	"fmt"
	"strings"
)

func addPortToURL(url string, port string) string {
	if port == "" {
		return url
	}

	const numSlashes = 3
	slashCount := 0
	insertIndex := len(url)

	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		// Find the position of the third '/'
		for i, ch := range url {
			if ch == '/' {
				slashCount++
				if slashCount == numSlashes {
					insertIndex = i
					break
				}
			}
		}

		// Construct the new URL
		if slashCount >= numSlashes {
			return url[:insertIndex] + ":" + port + url[insertIndex:]
		} else {
			return url + ":" + port
		}
	} else {
		// Find the position of the first '/'
		for i, ch := range url {
			if ch == '/' {
				insertIndex = i
				break
			}
		}

		// If no '/', add port at the end, otherwise add before the first '/'
		if insertIndex == len(url) {
			return url + ":" + port
		} else {
			return url[:insertIndex] + ":" + port + url[insertIndex:]
		}
	}
}

func main() {
	urls := []string{
		"https://costcenter.console.edgenext.com/entry.json",
		"http://example.com",
		"example.com/path",
		"example.io",
		"example.io/aaa/bbb/ccc",
		"example.cn",
		"example/aaa/bbb/ccc",
		"127.0.0.1",
		"127.0.0.1/aaa/bbb/ccc",
		"https://costcenter.console.edgenext.io/entry.json",
		"http://costcenter.console.edgenext.top/entry.json",
		"http://127.0.0.1/entry.json",
	}

	port := "8080"

	for _, url := range urls {
		newURL := addPortToURL(url, port)
		fmt.Println(newURL)
	}
}
