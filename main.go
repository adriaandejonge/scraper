package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

const path string = "./output"

func getEl(c *colly.Collector, el string, attr string) {
	c.OnHTML(el+"["+attr+"]", func(e *colly.HTMLElement) {
		link := e.Attr(attr)

		// Only relative URLs; no domain switches
		// Note that this way hardcoded domains get ignored
		if len(link) > 1 && link[:1] == "/" {
			e.Request.Visit(link)
		}

	})
}

func storeContent(body []byte, file string) {
	fullPath := path + file

	if fullPath[len(fullPath)-1:] == "/" || file == "" {
		fullPath = fullPath + "/index.html"
	}

	pos := strings.LastIndex(fullPath, "/")
	folderPath := fullPath[:pos]

	os.MkdirAll(folderPath, os.ModePerm)

	err := ioutil.WriteFile(fullPath, body, 0644)
	if err != nil {
		fmt.Println("Error writing file " + fullPath)
	} else {
		fmt.Println("Stored " + fullPath)
	}
}

func main() {
	c := colly.NewCollector(
		colly.MaxDepth(2),
	)

	c.Limit(&colly.LimitRule{
		Parallelism: 1,
		Delay:       500 * time.Millisecond,
	})

	getEl(c, "a", "href")
	getEl(c, "img", "src")
	getEl(c, "script", "src")
	getEl(c, "link", "href")

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("VISITING ", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		if r.StatusCode == 200 {
			file := r.Request.URL.EscapedPath()

			storeContent(r.Body, file)
		}
	})

	if len(os.Args) <= 1 {
		fmt.Println("Please provide the URL as a command line argument")
	} else {
		// To do: Some additional input validation would help
		c.Visit(os.Args[1])
	}
}
