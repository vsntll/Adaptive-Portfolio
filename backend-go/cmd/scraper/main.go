package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"linkedin-scraper/internal/config"
	"linkedin-scraper/internal/export"
	"linkedin-scraper/internal/models"
	"linkedin-scraper/internal/scraper"
)

func main() {
	var (
		profileURL = flag.String("profile", "", "LinkedIn profile URL to scrape")
		outputPath = flag.String("output", "data/output/profiles.csv", "Output CSV file path")
		configPath = flag.String("config", "configs/config.yaml", "Configuration file path")
		format     = flag.String("format", "csv", "Output format: csv, json, or both")
	)
	flag.Parse()

	if *profileURL == "" {
		log.Fatal("Profile URL is required. Use -profile flag to specify LinkedIn profile URL")
	}

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create scraper client
	client, err := scraper.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create scraper client: %v", err)
	}
	defer client.Close()

	// Login to LinkedIn
	log.Println("Logging in to LinkedIn...")
	err = client.Login()
	if err != nil {
		log.Fatalf("Failed to login to LinkedIn: %v", err)
	}

	// Scrape profile
	log.Printf("Scraping profile: %s", *profileURL)
	profile, err := client.ScrapeProfile(*profileURL)
	if err != nil {
		log.Fatalf("Failed to scrape profile: %v", err)
	}

	// Export data
	profiles := []*models.Profile{profile}
	
	// Create output directory if it doesn't exist
	if err := os.MkdirAll("data/output", 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	switch *format {
	case "csv":
		err = export.ToCSV(profiles, *outputPath)
	case "json":
		jsonPath := changeExtension(*outputPath, ".json")
		err = export.ToJSON(profiles, jsonPath)
	case "both":
		err = export.ToCSV(profiles, *outputPath)
		if err == nil {
			jsonPath := changeExtension(*outputPath, ".json")
			err = export.ToJSON(profiles, jsonPath)
		}
	default:
		log.Fatalf("Invalid format: %s. Use csv, json, or both", *format)
	}

	if err != nil {
		log.Fatalf("Failed to export data: %v", err)
	}

	fmt.Printf("Successfully scraped profile and saved to %s\n", *outputPath)
	fmt.Printf("Profile: %s\n", profile.Name)
	fmt.Printf("Headline: %s\n", profile.Headline)
	fmt.Printf("Location: %s\n", profile.Location)
}

func changeExtension(filename, newExt string) string {
	// Remove existing extension and add new one
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '.' {
			return filename[:i] + newExt
		}
	}
	return filename + newExt
}