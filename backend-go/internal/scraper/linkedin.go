package scraper

import (
	"fmt"
	"log"
	"strings"
	"time"

	"linkedin-scraper/internal/models"

	"github.com/tebeka/selenium"
)

// ScrapeProfile scrapes a LinkedIn profile and returns profile data
func (c *Client) ScrapeProfile(profileURL string) (*models.Profile, error) {
	log.Printf("Navigating to profile: %s", profileURL)
	
	// Navigate to profile
	err := c.driver.Get(profileURL)
	if err != nil {
		return nil, fmt.Errorf("failed to navigate to profile: %w", err)
	}

	// Wait for page to load
	time.Sleep(c.config.Delays.PageLoad)

	// Create new profile
	profile := models.NewProfile()
	profile.ProfileURL = profileURL

	// Perform human-like scrolling to load dynamic content
	err = c.humanLikeScrolling()
	if err != nil {
		log.Printf("Warning: Failed to scroll page: %v", err)
	}

	// Extract basic profile information
	err = c.extractBasicInfo(profile)
	if err != nil {
		log.Printf("Warning: Failed to extract basic info: %v", err)
	}

	c.randomDelay()

	// Extract experience
	err = c.extractExperience(profile)
	if err != nil {
		log.Printf("Warning: Failed to extract experience: %v", err)
	}

	c.randomDelay()

	// Extract education
	err = c.extractEducation(profile)
	if err != nil {
		log.Printf("Warning: Failed to extract education: %v", err)
	}

	c.randomDelay()

	// Extract skills
	err = c.extractSkills(profile)
	if err != nil {
		log.Printf("Warning: Failed to extract skills: %v", err)
	}

	// Validate profile has minimum required data
	if !profile.Validate() {
		return nil, fmt.Errorf("profile validation failed - insufficient data extracted")
	}

	log.Printf("Successfully scraped profile for: %s", profile.Name)
	return profile, nil
}

// extractBasicInfo extracts name, headline, location, and about section
func (c *Client) extractBasicInfo(profile *models.Profile) error {
	// Extract name
	nameSelectors := []string{
		"h1.text-heading-xlarge",
		"h1[data-generated-suggestion-target]",
		".pv-text-details__left-panel h1",
		".ph5 h1",
	}
	
	name, err := c.extractTextBySelectors(nameSelectors)
	if err == nil && name != "" {
		profile.Name = strings.TrimSpace(name)
	}

	// Extract headline
	headlineSelectors := []string{
		".text-body-medium.break-words",
		".pv-text-details__left-panel .text-body-medium",
		".ph5 .text-body-medium",
		"[data-generated-suggestion-target] + div",
	}
	
	headline, err := c.extractTextBySelectors(headlineSelectors)
	if err == nil && headline != "" {
		profile.Headline = strings.TrimSpace(headline)
	}

	// Extract location
	locationSelectors := []string{
		".text-body-small.inline.t-black--light.break-words",
		".pv-text-details__left-panel .text-body-small",
		".ph5 .text-body-small",
		"span.text-body-small.inline",
	}
	
	location, err := c.extractTextBySelectors(locationSelectors)
	if err == nil && location != "" {
		profile.Location = strings.TrimSpace(location)
	}

	// Extract about section
	aboutSelectors := []string{
		"#about ~ * .inline-show-more-text span[aria-hidden='true']",
		".pv-shared-text-with-see-more span[aria-hidden='true']",
		"section[data-section='summary'] .pv-shared-text-with-see-more",
		".core-section-container__content .inline-show-more-text",
	}
	
	about, err := c.extractTextBySelectors(aboutSelectors)
	if err == nil && about != "" {
		profile.About = strings.TrimSpace(about)
	}

	return nil
}

// extractExperience extracts work experience
func (c *Client) extractExperience(profile *models.Profile) error {
	// Try to find experience section
	experienceSelectors := []string{
		"#experience ~ * .pvs-list__item--line-separated",
		".experience-section .pv-entity__summary-info",
		"section[data-section='experience'] .pv-entity__summary-info",
	}

	var experienceElements []selenium.WebElement
	var err error

	for _, selector := range experienceSelectors {
		elements, findErr := c.driver.FindElements(selenium.ByCSSSelector, selector)
		if findErr == nil && len(elements) > 0 {
			experienceElements = elements
			break
		}
	}

	if len(experienceElements) == 0 {
		// Try alternative XPath approach
		elements, err := c.driver.FindElements(selenium.ByXPATH, "//section[contains(@id, 'experience')]//li[contains(@class, 'artdeco-list__item')]")
		if err != nil || len(elements) == 0 {
			return fmt.Errorf("no experience elements found")
		}
		experienceElements = elements
	}

	// Extract experience data
	for i, element := range experienceElements {
		if i >= 5 { // Limit to first 5 experiences
			break
		}

		experience := models.Experience{}

		// Extract title
		titleElement, err := element.FindElement(selenium.ByCSSSelector, ".mr1.t-bold span[aria-hidden='true']")
		if err == nil {
			title, _ := titleElement.Text()
			experience.Title = strings.TrimSpace(title)
		}

		// Extract company
		companyElement, err := element.FindElement(selenium.ByCSSSelector, ".t-14.t-normal span[aria-hidden='true']")
		if err == nil {
			company, _ := companyElement.Text()
			experience.Company = strings.TrimSpace(company)
		}

		// Extract duration
		durationElement, err := element.FindElement(selenium.ByCSSSelector, ".t-12.t-normal--light span[aria-hidden='true']")
		if err == nil {
			duration, _ := durationElement.Text()
			experience.Duration = strings.TrimSpace(duration)
		}

		if experience.Title != "" || experience.Company != "" {
			profile.AddExperience(experience)
		}
	}

	return nil
}

// extractEducation extracts educational background
func (c *Client) extractEducation(profile *models.Profile) error {
	// Try to find education section
	educationSelectors := []string{
		"#education ~ * .pvs-list__item--line-separated",
		".education-section .pv-entity__summary-info",
		"section[data-section='education'] .pv-entity__summary-info",
	}

	var educationElements []selenium.WebElement

	for _, selector := range educationSelectors {
		elements, err := c.driver.FindElements(selenium.ByCSSSelector, selector)
		if err == nil && len(elements) > 0 {
			educationElements = elements
			break
		}
	}

	if len(educationElements) == 0 {
		// Try alternative XPath approach
		elements, err := c.driver.FindElements(selenium.ByXPATH, "//section[contains(@id, 'education')]//li[contains(@class, 'artdeco-list__item')]")
		if err != nil || len(elements) == 0 {
			return fmt.Errorf("no education elements found")
		}
		educationElements = elements
	}

	// Extract education data
	for i, element := range educationElements {
		if i >= 3 { // Limit to first 3 education entries
			break
		}

		education := models.Education{}

		// Extract school name
		schoolElement, err := element.FindElement(selenium.ByCSSSelector, ".mr1.t-bold span[aria-hidden='true']")
		if err == nil {
			school, _ := schoolElement.Text()
			education.School = strings.TrimSpace(school)
		}

		// Extract degree
		degreeElement, err := element.FindElement(selenium.ByCSSSelector, ".t-14.t-normal span[aria-hidden='true']")
		if err == nil {
			degree, _ := degreeElement.Text()
			education.Degree = strings.TrimSpace(degree)
		}

		// Extract duration
		durationElement, err := element.FindElement(selenium.ByCSSSelector, ".t-12.t-normal--light span[aria-hidden='true']")
		if err == nil {
			duration, _ := durationElement.Text()
			education.Duration = strings.TrimSpace(duration)
		}

		if education.School != "" || education.Degree != "" {
			profile.AddEducation(education)
		}
	}

	return nil
}

// extractSkills extracts skills from the profile
func (c *Client) extractSkills(profile *models.Profile) error {
	// Try to find skills section
	skillSelectors := []string{
		"#skills ~ * .mr1.t-bold span[aria-hidden='true']",
		".skills-section .pv-skill-category-entity__name span",
		"section[data-section='skills'] .pv-skill-category-entity__name span",
	}

	var skillElements []selenium.WebElement

	for _, selector := range skillSelectors {
		elements, err := c.driver.FindElements(selenium.ByCSSSelector, selector)
		if err == nil && len(elements) > 0 {
			skillElements = elements
			break
		}
	}

	if len(skillElements) == 0 {
		// Try alternative XPath approach
		elements, err := c.driver.FindElements(selenium.ByXPATH, "//section[contains(@id, 'skills')]//span[contains(@class, 't-bold')]")
		if err != nil || len(elements) == 0 {
			return fmt.Errorf("no skill elements found")
		}
		skillElements = elements
	}

	// Extract skills
	for i, element := range skillElements {
		if i >= 10 { // Limit to first 10 skills
			break
		}

		skill, err := element.Text()
		if err == nil && skill != "" {
			profile.AddSkill(skill)
		}
	}

	return nil
}

// extractTextBySelectors tries multiple CSS selectors to extract text
func (c *Client) extractTextBySelectors(selectors []string) (string, error) {
	for _, selector := range selectors {
		element, err := c.driver.FindElement(selenium.ByCSSSelector, selector)
		if err == nil {
			text, textErr := element.Text()
			if textErr == nil && strings.TrimSpace(text) != "" {
				return text, nil
			}
		}
	}
	return "", fmt.Errorf("text not found with any selector")
}