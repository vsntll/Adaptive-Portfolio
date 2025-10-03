package export

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"linkedin-scraper/internal/models"
)

// ToCSV exports profiles to CSV format
func ToCSV(profiles []*models.Profile, filename string) error {
	// Ensure directory exists
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create CSV file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	header := []string{
		"Name",
		"Headline", 
		"Location",
		"About",
		"Experience",
		"Education",
		"Skills",
		"Profile URL",
		"Scraped At",
	}
	
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write profile data
	for _, profile := range profiles {
		record := []string{
			profile.Name,
			profile.Headline,
			profile.Location,
			profile.About,
			profile.GetExperienceAsString(),
			profile.GetEducationAsString(),
			profile.GetSkillsAsString(),
			profile.ProfileURL,
			profile.ScrapedAt.Format("2006-01-02 15:04:05"),
		}
		
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	return nil
}

// ToJSON exports profiles to JSON format
func ToJSON(profiles []*models.Profile, filename string) error {
	// Ensure directory exists
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create JSON file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create JSON file: %w", err)
	}
	defer file.Close()

	// Encode profiles as JSON
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty print JSON

	// Wrap profiles in a container object
	data := map[string]interface{}{
		"profiles": profiles,
		"count":    len(profiles),
		"exported_at": profiles[0].ScrapedAt.Format("2006-01-02 15:04:05"),
	}

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

// ToCSVDetailed exports profiles to CSV with detailed experience and education
func ToCSVDetailed(profiles []*models.Profile, filename string) error {
	// Ensure directory exists
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create CSV file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write detailed CSV header
	header := []string{
		"Name",
		"Headline", 
		"Location",
		"About",
		"Experience Title",
		"Experience Company",
		"Experience Duration",
		"Experience Location",
		"Education School",
		"Education Degree",
		"Education Duration",
		"Skills",
		"Profile URL",
		"Scraped At",
	}
	
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write profile data (one row per experience/education combination)
	for _, profile := range profiles {
		maxRows := max(len(profile.Experience), len(profile.Education))
		if maxRows == 0 {
			maxRows = 1 // At least one row per profile
		}

		for i := 0; i < maxRows; i++ {
			record := []string{
				profile.Name,
				profile.Headline,
				profile.Location,
				profile.About,
			}

			// Add experience data if available
			if i < len(profile.Experience) {
				exp := profile.Experience[i]
				record = append(record, 
					exp.Title,
					exp.Company,
					exp.Duration,
					exp.Location,
				)
			} else {
				record = append(record, "", "", "", "")
			}

			// Add education data if available
			if i < len(profile.Education) {
				edu := profile.Education[i]
				record = append(record,
					edu.School,
					edu.Degree,
					edu.Duration,
				)
			} else {
				record = append(record, "", "", "")
			}

			// Add skills and metadata (only on first row)
			if i == 0 {
				record = append(record,
					profile.GetSkillsAsString(),
					profile.ProfileURL,
					profile.ScrapedAt.Format("2006-01-02 15:04:05"),
				)
			} else {
				record = append(record, "", "", "")
			}
			
			if err := writer.Write(record); err != nil {
				return fmt.Errorf("failed to write CSV record: %w", err)
			}
		}
	}

	return nil
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ExportSummary creates a summary report of the exported data
func ExportSummary(profiles []*models.Profile, filename string) error {
	// Ensure directory exists
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create summary file: %w", err)
	}
	defer file.Close()

	// Calculate statistics
	totalProfiles := len(profiles)
	totalExperiences := 0
	totalEducation := 0
	totalSkills := 0

	for _, profile := range profiles {
		totalExperiences += len(profile.Experience)
		totalEducation += len(profile.Education)
		totalSkills += len(profile.Skills)
	}

	avgExperiences := float64(totalExperiences) / float64(totalProfiles)
	avgEducation := float64(totalEducation) / float64(totalProfiles)
	avgSkills := float64(totalSkills) / float64(totalProfiles)

	// Write summary
	summary := fmt.Sprintf(`LinkedIn Profile Scraping Summary
=====================================

Total Profiles: %d
Total Experience Entries: %d
Total Education Entries: %d
Total Skills: %d

Average Experience per Profile: %.2f
Average Education per Profile: %.2f
Average Skills per Profile: %.2f

Profiles:
`, totalProfiles, totalExperiences, totalEducation, totalSkills, avgExperiences, avgEducation, avgSkills)

	for i, profile := range profiles {
		summary += fmt.Sprintf("%d. %s (%s)\n", i+1, profile.Name, profile.ProfileURL)
	}

	summary += fmt.Sprintf("\nGenerated at: %s\n", profiles[0].ScrapedAt.Format("2006-01-02 15:04:05"))

	_, err = file.WriteString(summary)
	return err
}