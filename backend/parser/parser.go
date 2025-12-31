package parser

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/ledongthuc/pdf"
)

// ResumeData defines the structure of our output JSON
type ResumeData struct {
	Name       string   `json:"name"`
	Email      string   `json:"email"`
	Summary    string   `json:"summary"`
	Skills     []string `json:"skills"`
	Experience []string `json:"experience"`
	Education  []string `json:"education"`
}

// ParsePDF reads the PDF and returns structured data
func ParsePDF(path string) (ResumeData, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return ResumeData{}, err
	}
	defer f.Close()

	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return ResumeData{}, err
	}
	buf.ReadFrom(b)
	text := buf.String()

	return heuristicParse(text), nil
}

// heuristicParse splits text based on common resume headers
func heuristicParse(text string) ResumeData {
	data := ResumeData{}
	lines := strings.Split(text, "\n")

	// Regex for email
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)

	// Section keywords
	sectionKeywords := map[string]string{
		"experience":   "experience",
		"work history": "experience",
		"education":    "education",
		"skills":       "skills",
		"technologies": "skills",
		"summary":      "summary",
		"profile":      "summary",
	}

	var currentSection string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Section detection
		isSectionHeader := false
		for keyword, section := range sectionKeywords {
			if strings.Contains(strings.ToLower(line), keyword) {
				currentSection = section
				isSectionHeader = true
				break
			}
		}
		if isSectionHeader {
			continue
		}

		// Data extraction
		switch currentSection {
		case "summary":
			data.Summary += line + " "
		case "skills":
			skills := strings.Split(line, ",")
			for _, skill := range skills {
				data.Skills = append(data.Skills, strings.TrimSpace(skill))
			}
		case "experience":
			data.Experience = append(data.Experience, line)
		case "education":
			data.Education = append(data.Education, line)
		default:
			// Header section
			if email := emailRegex.FindString(line); email != "" {
				data.Email = email
				// Assume the rest of the line is part of the name
				namePart := strings.TrimSpace(strings.Replace(line, email, "", 1))
				if namePart != "" {
					data.Name = strings.TrimSpace(data.Name + " " + namePart)
				}
			} else if data.Name == "" {
				data.Name = line
			}
		}
	}

	// Clean up summary
	data.Summary = strings.TrimSpace(data.Summary)

	return data
}