package models

import (
	"strings"
	"time"
)

// Profile represents a LinkedIn profile
type Profile struct {
	Name       string       `json:"name" csv:"name"`
	Headline   string       `json:"headline" csv:"headline"`
	Location   string       `json:"location" csv:"location"`
	About      string       `json:"about" csv:"about"`
	Experience []Experience `json:"experience" csv:"experience"`
	Education  []Education  `json:"education" csv:"education"`
	Skills     []string     `json:"skills" csv:"skills"`
	ScrapedAt  time.Time    `json:"scraped_at" csv:"scraped_at"`
	ProfileURL string       `json:"profile_url" csv:"profile_url"`
}

// Experience represents work experience
type Experience struct {
	Title       string `json:"title" csv:"title"`
	Company     string `json:"company" csv:"company"`
	Duration    string `json:"duration" csv:"duration"`
	Location    string `json:"location" csv:"location"`
	Description string `json:"description" csv:"description"`
}

// Education represents educational background
type Education struct {
	School      string `json:"school" csv:"school"`
	Degree      string `json:"degree" csv:"degree"`
	Duration    string `json:"duration" csv:"duration"`
	Description string `json:"description" csv:"description"`
}

// NewProfile creates a new profile instance
func NewProfile() *Profile {
	return &Profile{
		Experience: make([]Experience, 0),
		Education:  make([]Education, 0),
		Skills:     make([]string, 0),
		ScrapedAt:  time.Now(),
	}
}

// AddExperience adds work experience to the profile
func (p *Profile) AddExperience(exp Experience) {
	p.Experience = append(p.Experience, exp)
}

// AddEducation adds education to the profile
func (p *Profile) AddEducation(edu Education) {
	p.Education = append(p.Education, edu)
}

// AddSkill adds a skill to the profile
func (p *Profile) AddSkill(skill string) {
	if skill != "" && !p.hasSkill(skill) {
		p.Skills = append(p.Skills, strings.TrimSpace(skill))
	}
}

// hasSkill checks if a skill already exists
func (p *Profile) hasSkill(skill string) bool {
	for _, s := range p.Skills {
		if strings.EqualFold(s, skill) {
			return true
		}
	}
	return false
}

// GetExperienceAsString returns experience as a formatted string for CSV
func (p *Profile) GetExperienceAsString() string {
	if len(p.Experience) == 0 {
		return ""
	}
	
	var parts []string
	for _, exp := range p.Experience {
		part := exp.Title + " at " + exp.Company
		if exp.Duration != "" {
			part += " (" + exp.Duration + ")"
		}
		parts = append(parts, part)
	}
	return strings.Join(parts, "; ")
}

// GetEducationAsString returns education as a formatted string for CSV
func (p *Profile) GetEducationAsString() string {
	if len(p.Education) == 0 {
		return ""
	}
	
	var parts []string
	for _, edu := range p.Education {
		part := edu.Degree + " from " + edu.School
		if edu.Duration != "" {
			part += " (" + edu.Duration + ")"
		}
		parts = append(parts, part)
	}
	return strings.Join(parts, "; ")
}

// GetSkillsAsString returns skills as a comma-separated string
func (p *Profile) GetSkillsAsString() string {
	return strings.Join(p.Skills, ", ")
}

// Validate checks if the profile has minimum required data
func (p *Profile) Validate() bool {
	return p.Name != "" && p.ProfileURL != ""
}