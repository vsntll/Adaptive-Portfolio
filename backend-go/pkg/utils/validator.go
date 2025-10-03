package utils

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// ValidateLinkedInURL validates if a URL is a valid LinkedIn profile URL
func ValidateLinkedInURL(profileURL string) error {
	if profileURL == "" {
		return fmt.Errorf("profile URL cannot be empty")
	}

	// Parse URL
	parsedURL, err := url.Parse(profileURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	// Check if it's a LinkedIn URL
	if parsedURL.Host != "www.linkedin.com" && parsedURL.Host != "linkedin.com" {
		return fmt.Errorf("URL must be from linkedin.com domain")
	}

	// Check if it's a profile URL
	profilePattern := regexp.MustCompile(`^/in/[a-zA-Z0-9\-]+/?$`)
	if !profilePattern.MatchString(parsedURL.Path) {
		return fmt.Errorf("URL must be a LinkedIn profile URL (format: /in/username)")
	}

	return nil
}

// NormalizeLinkedInURL normalizes a LinkedIn profile URL
func NormalizeLinkedInURL(profileURL string) (string, error) {
	// Add https:// if missing
	if !strings.HasPrefix(profileURL, "http://") && !strings.HasPrefix(profileURL, "https://") {
		profileURL = "https://" + profileURL
	}

	// Parse URL
	parsedURL, err := url.Parse(profileURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL format: %w", err)
	}

	// Normalize host
	if parsedURL.Host == "linkedin.com" {
		parsedURL.Host = "www.linkedin.com"
	}

	// Ensure HTTPS
	parsedURL.Scheme = "https"

	// Remove trailing slash if present
	parsedURL.Path = strings.TrimSuffix(parsedURL.Path, "/")

	// Remove query parameters and fragments
	parsedURL.RawQuery = ""
	parsedURL.Fragment = ""

	return parsedURL.String(), nil
}

// ValidateProfileData validates that profile data meets minimum requirements
func ValidateProfileData(name, profileURL string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("profile name is required")
	}

	if strings.TrimSpace(profileURL) == "" {
		return fmt.Errorf("profile URL is required")
	}

	return ValidateLinkedInURL(profileURL)
}

// SanitizeText cleans and sanitizes text data
func SanitizeText(text string) string {
	// Remove excessive whitespace
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
	
	// Trim whitespace
	text = strings.TrimSpace(text)
	
	// Remove common LinkedIn artifacts
	text = strings.ReplaceAll(text, "…", "...")
	text = strings.ReplaceAll(text, """, "\"")
	text = strings.ReplaceAll(text, """, "\"")
	text = strings.ReplaceAll(text, "'", "'")
	text = strings.ReplaceAll(text, "'", "'")
	
	return text
}

// ValidateEmail validates email format
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailPattern.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

// ValidatePassword validates password strength
func ValidatePassword(password string) error {
	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	if len(password) < 6 {
		return fmt.Errorf("password must be at least 6 characters long")
	}

	return nil
}

// ExtractUsernameFromURL extracts the username from a LinkedIn profile URL
func ExtractUsernameFromURL(profileURL string) (string, error) {
	parsedURL, err := url.Parse(profileURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL format: %w", err)
	}

	// Extract username from path
	pathParts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
	if len(pathParts) < 2 || pathParts[0] != "in" {
		return "", fmt.Errorf("invalid LinkedIn profile URL format")
	}

	return pathParts[1], nil
}

// IsEmptyOrWhitespace checks if a string is empty or contains only whitespace
func IsEmptyOrWhitespace(str string) bool {
	return strings.TrimSpace(str) == ""
}

// TruncateText truncates text to a maximum length with ellipsis
func TruncateText(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}
	
	if maxLength <= 3 {
		return text[:maxLength]
	}
	
	return text[:maxLength-3] + "..."
}

// ContainsAny checks if a string contains any of the specified substrings
func ContainsAny(text string, substrings []string) bool {
	lowerText := strings.ToLower(text)
	for _, substring := range substrings {
		if strings.Contains(lowerText, strings.ToLower(substring)) {
			return true
		}
	}
	return false
}

// RemoveDuplicates removes duplicate strings from a slice while preserving order
func RemoveDuplicates(slice []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(slice))
	
	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	
	return result
}