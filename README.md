# scraper
Based on Colly - store site as static files

To run in interpreted mode:
```
go get -u github.com/gocolly/colly/...
go run main.go https://your-url-here/
```

To compile:
```
go get -u github.com/gocolly/colly/...
go build main.go
./main https://your-url-here/
```

Note: this is PoC only hence a bit rough around the edges. As JavaScript is not executed during crawling, some links/images may be missing.
