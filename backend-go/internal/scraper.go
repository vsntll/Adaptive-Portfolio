package scraper

import (
    "github.com/gocolly/colly/v2"
)

// ScrapeLinks scrapes all anchor href links from the given URL.
func ScrapeLinks(url string) ([]string, error) {
    var links []string

    c := colly.NewCollector()

    c.OnHTML("a[href]", func(e *colly.HTMLElement) {
        links = append(links, e.Attr("href"))
    })

    err := c.Visit(url)
    if err != nil {
        return nil, err
    }

    return links, nil
}
