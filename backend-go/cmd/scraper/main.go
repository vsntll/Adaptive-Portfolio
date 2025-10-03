package main

import (
    "fmt"
    "log"

    "adaptive-portfolio/internal/scraper"
)

func main() {
    urls := []string{
        "https://example.com",
        // Add other URLs you want to scrape here
    }

    for _, url := range urls {
        data, err := scraper.ScrapeLinks(url)
        if err != nil {
            log.Printf("Error scraping %s: %v\n", url, err)
            continue
        }
        fmt.Printf("Links scraped from %s:\n", url)
        for _, link := range data {
            fmt.Println(link)
        }
    }
}
