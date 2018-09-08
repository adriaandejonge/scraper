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

		fmt.Println("   "+el, link)

		if len(link) > 1 && link[:1] == "/" {
			e.Request.Visit(link)
		}

	})
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func storeContent(body []byte, file string) {
	fullPath := path + file
	pos := strings.LastIndex(fullPath, "/")
	folderPath := fullPath[:pos]

	os.MkdirAll(folderPath, os.ModePerm)

	fi, err := os.Lstat(fullPath)
	//check(err)

	if err != os.ErrNotExist && fi.Mode().IsDir() {
		// ASSUMPTION: .html
		fullPath = fullPath + "/index.html"
	}

	err = ioutil.WriteFile(fullPath, body, 0644)
	check(err)
}

func main() {
	c := colly.NewCollector(
		colly.MaxDepth(2),
	)

	c.Limit(&colly.LimitRule{
		Parallelism: 2,
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

			// if file[len(file)-1:] == "/" {
			// 	file = file + "index.html"
			// }
			storeContent(r.Body, file)
			fmt.Println("response received", file)
		}

	})

	c.Visit("https://www.singtel.com/")
}
