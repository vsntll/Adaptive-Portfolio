package scraper

import (
	"math/rand"
	"time"

	"github.com/tebeka/selenium"
)

// randomDelay implements anti-detection random delays
func (c *Client) randomDelay() {
	min := c.config.Delays.MinDelay
	max := c.config.Delays.MaxDelay
	
	// Generate random delay between min and max
	delayRange := int64(max - min)
	if delayRange <= 0 {
		time.Sleep(min)
		return
	}
	
	randomDelay := time.Duration(rand.Int63n(delayRange)) + min
	time.Sleep(randomDelay)
}

// shortDelay implements shorter delays for quick actions
func (c *Client) shortDelay() {
	delay := time.Duration(rand.Intn(1000)+500) * time.Millisecond
	time.Sleep(delay)
}

// humanLikeScrolling simulates human scrolling behavior
func (c *Client) humanLikeScrolling() error {
	// Get page height first
	totalHeight, err := c.driver.ExecuteScript("return document.body.scrollHeight", nil)
	if err != nil {
		return err
	}
	
	height, ok := totalHeight.(int64)
	if !ok {
		height = 3000 // fallback height
	}
	
	// Scroll in chunks with random delays
	scrollStep := int64(300)
	currentPosition := int64(0)
	
	for currentPosition < height {
		// Random scroll amount
		scrollAmount := scrollStep + int64(rand.Intn(200)-100)
		if scrollAmount < 100 {
			scrollAmount = 100
		}
		
		script := "window.scrollBy(0, arguments[0]);"
		_, err := c.driver.ExecuteScript(script, []interface{}{scrollAmount})
		if err != nil {
			return err
		}
		
		currentPosition += scrollAmount
		
		// Random delay between scrolls
		scrollDelay := c.config.Delays.ScrollDelay + time.Duration(rand.Intn(500))*time.Millisecond
		time.Sleep(scrollDelay)
		
		// Occasionally pause longer (simulate reading)
		if rand.Intn(5) == 0 {
			time.Sleep(time.Duration(rand.Intn(2000)+1000) * time.Millisecond)
		}
	}
	
	// Scroll back to top
	_, err = c.driver.ExecuteScript("window.scrollTo(0, 0);", nil)
	return err
}

// rotateUserAgent rotates user agent to avoid detection
func (c *Client) rotateUserAgent() error {
	if len(c.config.UserAgents) == 0 {
		return nil
	}
	
	// Select random user agent
	userAgent := c.config.UserAgents[rand.Intn(len(c.config.UserAgents))]
	
	// Set user agent via CDP (Chrome DevTools Protocol)
	script := `
		Object.defineProperty(navigator, 'userAgent', {
			get: function() { return arguments[0]; }
		});
	`
	
	_, err := c.driver.ExecuteScript(script, []interface{}{userAgent})
	return err
}

// simulateMouseMovement simulates random mouse movements
func (c *Client) simulateMouseMovement() error {
	// Get viewport dimensions
	width, err := c.driver.ExecuteScript("return window.innerWidth", nil)
	if err != nil {
		return err
	}
	
	height, err := c.driver.ExecuteScript("return window.innerHeight", nil)
	if err != nil {
		return err
	}
	
	w, ok1 := width.(int64)
	h, ok2 := height.(int64)
	
	if !ok1 || !ok2 {
		w, h = 1920, 1080 // fallback dimensions
	}
	
	// Generate random mouse positions
	for i := 0; i < 3; i++ {
		x := rand.Int63n(w)
		y := rand.Int63n(h)
		
		// Move mouse to random position
		script := `
			var event = new MouseEvent('mousemove', {
				clientX: arguments[0],
				clientY: arguments[1],
				bubbles: true
			});
			document.dispatchEvent(event);
		`
		
		c.driver.ExecuteScript(script, []interface{}{x, y})
		time.Sleep(time.Duration(rand.Intn(500)+100) * time.Millisecond)
	}
	
	return nil
}

// waitWithRandomization waits with slight randomization
func (c *Client) waitWithRandomization(baseDelay time.Duration) {
	// Add random variance of ±20%
	variance := int64(float64(baseDelay) * 0.2)
	randomVariance := time.Duration(rand.Int63n(2*variance) - variance)
	finalDelay := baseDelay + randomVariance
	
	if finalDelay < 0 {
		finalDelay = baseDelay
	}
	
	time.Sleep(finalDelay)
}

// clickWithDelay clicks an element with human-like delay
func (c *Client) clickWithDelay(element selenium.WebElement) error {
	// Simulate mouse movement before click
	c.simulateMouseMovement()
	
	// Short delay before click
	c.shortDelay()
	
	// Click element
	err := element.Click()
	if err != nil {
		return err
	}
	
	// Short delay after click
	c.shortDelay()
	
	return nil
}

// typeWithDelay types text with human-like delays
func (c *Client) typeWithDelay(element selenium.WebElement, text string) error {
	// Clear field first
	err := element.Clear()
	if err != nil {
		return err
	}
	
	c.shortDelay()
	
	// Type character by character with random delays
	for _, char := range text {
		err = element.SendKeys(string(char))
		if err != nil {
			return err
		}
		
		// Random delay between keystrokes (20-100ms)
		delay := time.Duration(rand.Intn(80)+20) * time.Millisecond
		time.Sleep(delay)
	}
	
	c.shortDelay()
	return nil
}

// antiDetectionSetup configures additional anti-detection measures
func (c *Client) antiDetectionSetup() error {
	// Disable automation indicator
	script := `
		Object.defineProperty(navigator, 'webdriver', {
			get: () => undefined,
		});
		
		// Remove automation flag from Chrome
		window.chrome = {
			runtime: {},
		};
		
		// Override permissions API
		Object.defineProperty(navigator, 'permissions', {
			get: () => ({
				query: () => Promise.resolve({ state: 'granted' }),
			}),
		});
		
		// Override plugins length
		Object.defineProperty(navigator, 'plugins', {
			get: () => [1, 2, 3, 4, 5],
		});
		
		// Override languages
		Object.defineProperty(navigator, 'languages', {
			get: () => ['en-US', 'en'],
		});
	`
	
	_, err := c.driver.ExecuteScript(script, nil)
	return err
}

// checkForCaptcha checks if CAPTCHA or verification is present
func (c *Client) checkForCaptcha() (bool, error) {
	captchaSelectors := []string{
		"iframe[src*='recaptcha']",
		"div[class*='captcha']",
		"div[class*='challenge']",
		"div[id*='captcha']",
		".challenge-page",
		"[data-test='challenge']",
	}
	
	for _, selector := range captchaSelectors {
		elements, err := c.driver.FindElements(selenium.ByCSSSelector, selector)
		if err == nil && len(elements) > 0 {
			return true, nil
		}
	}
	
	return false, nil
}

// handleRateLimit handles rate limiting scenarios
func (c *Client) handleRateLimit() error {
	// Check for rate limit indicators
	rateLimitSelectors := []string{
		"[data-test='rate-limit']",
		".rate-limit",
		"div:contains('too many requests')",
		"div:contains('Please wait')",
	}
	
	for _, selector := range rateLimitSelectors {
		elements, err := c.driver.FindElements(selenium.ByCSSSelector, selector)
		if err == nil && len(elements) > 0 {
			// Implement exponential backoff
			backoffTime := time.Duration(rand.Intn(300)+300) * time.Second
			time.Sleep(backoffTime)
			return nil
		}
	}
	
	return nil
}