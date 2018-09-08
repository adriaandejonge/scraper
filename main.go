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
	pos := strings.LastIndex(fullPath, "/")
	folderPath := fullPath[:pos]

	os.MkdirAll(folderPath, os.ModePerm)

	if _, err := os.Stat(fullPath); err != nil {
		if !os.IsNotExist(err) {
			fullPath = fullPath + "/index.html"
		}
	}

	if fullPath[len(fullPath)-1:] == "/" {
		fullPath = fullPath + "/index.html"
	}

	err := ioutil.WriteFile(fullPath, body, 0644)
	if err != nil {
		fmt.Println("Error writing file " + fullPath)
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
			fmt.Println("response received", file)
		}
	})

	c.Visit("https://www.singtel.com/")
}
