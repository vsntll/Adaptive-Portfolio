package scraper

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"linkedin-scraper/internal/config"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

// Client represents the web scraper client
type Client struct {
	driver  selenium.WebDriver
	service *selenium.Service
	config  *config.Config
}

// NewClient creates a new scraper client
func NewClient(cfg *config.Config) (*Client, error) {
	// Setup Chrome options
	caps := selenium.Capabilities{"browserName": "chrome"}
	chromeCaps := chrome.Capabilities{
		Args: []string{
			"--no-sandbox",
			"--disable-blink-features=AutomationControlled",
			"--disable-gpu",
			"--disable-dev-shm-usage",
			"--disable-web-security",
			"--disable-features=VizDisplayCompositor",
			fmt.Sprintf("--window-size=%s", cfg.WebDriver.WindowSize),
		},
	}

	if cfg.WebDriver.Headless {
		chromeCaps.Args = append(chromeCaps.Args, "--headless")
	}

	caps.AddChrome(chromeCaps)

	// Start ChromeDriver service
	service, err := selenium.NewChromeDriverService(cfg.WebDriver.ChromeDriverPath, cfg.WebDriver.Port)
	if err != nil {
		return nil, fmt.Errorf("failed to start ChromeDriver service: %w", err)
	}

	// Connect to WebDriver
	driver, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", cfg.WebDriver.Port))
	if err != nil {
		service.Stop()
		return nil, fmt.Errorf("failed to connect to WebDriver: %w", err)
	}

	// Set timeouts
	driver.SetImplicitWaitTimeout(cfg.Delays.ElementSearch)
	driver.SetPageLoadTimeout(time.Duration(cfg.WebDriver.Timeout) * time.Second)

	client := &Client{
		driver:  driver,
		service: service,
		config:  cfg,
	}

	// Set initial user agent
	err = client.rotateUserAgent()
	if err != nil {
		log.Printf("Warning: Failed to set user agent: %v", err)
	}

	return client, nil
}

// Close cleanup resources
func (c *Client) Close() {
	if c.driver != nil {
		c.driver.Quit()
	}
	if c.service != nil {
		c.service.Stop()
	}
}

// Login authenticates with LinkedIn
func (c *Client) Login() error {
	// Navigate to LinkedIn login page
	err := c.driver.Get(c.config.LinkedIn.LoginURL)
	if err != nil {
		return fmt.Errorf("failed to navigate to login page: %w", err)
	}

	c.randomDelay()

	// Find and fill email field
	emailField, err := c.waitForElement(selenium.ByID, "username", c.config.Delays.ElementSearch)
	if err != nil {
		return fmt.Errorf("failed to find email field: %w", err)
	}

	err = emailField.SendKeys(c.config.LinkedIn.Email)
	if err != nil {
		return fmt.Errorf("failed to enter email: %w", err)
	}

	c.shortDelay()

	// Find and fill password field
	passwordField, err := c.waitForElement(selenium.ByID, "password", c.config.Delays.ElementSearch)
	if err != nil {
		return fmt.Errorf("failed to find password field: %w", err)
	}

	err = passwordField.SendKeys(c.config.LinkedIn.Password)
	if err != nil {
		return fmt.Errorf("failed to enter password: %w", err)
	}

	c.shortDelay()

	// Click login button
	loginButton, err := c.waitForElement(selenium.ByXPATH, "//button[@type='submit']", c.config.Delays.ElementSearch)
	if err != nil {
		return fmt.Errorf("failed to find login button: %w", err)
	}

	err = loginButton.Click()
	if err != nil {
		return fmt.Errorf("failed to click login button: %w", err)
	}

	// Wait for login to complete
	time.Sleep(c.config.Delays.PageLoad)

	// Check if login was successful by looking for the feed or profile
	_, err = c.waitForElement(selenium.ByXPATH, "//a[contains(@href, '/in/') or contains(@href, '/feed/')]", 10*time.Second)
	if err != nil {
		return fmt.Errorf("login appears to have failed or requires additional verification: %w", err)
	}

	log.Println("Successfully logged in to LinkedIn")
	return nil
}

// waitForElement waits for an element to be present and returns it
func (c *Client) waitForElement(by, value string, timeout time.Duration) (selenium.WebElement, error) {
	endTime := time.Now().Add(timeout)
	
	for time.Now().Before(endTime) {
		element, err := c.driver.FindElement(by, value)
		if err == nil {
			return element, nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	
	return nil, fmt.Errorf("element not found after %v: %s=%s", timeout, by, value)
}

// waitForElements waits for elements to be present and returns them
func (c *Client) waitForElements(by, value string, timeout time.Duration) ([]selenium.WebElement, error) {
	endTime := time.Now().Add(timeout)
	
	for time.Now().Before(endTime) {
		elements, err := c.driver.FindElements(by, value)
		if err == nil && len(elements) > 0 {
			return elements, nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	
	return nil, fmt.Errorf("elements not found after %v: %s=%s", timeout, by, value)
}

// parseWindowSize parses window size string and sets it
func parseWindowSize(size string) (int, int, error) {
	parts := strings.Split(size, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid window size format: %s", size)
	}

	width, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, fmt.Errorf("invalid width: %s", parts[0])
	}

	height, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, fmt.Errorf("invalid height: %s", parts[1])
	}

	return width, height, nil
}