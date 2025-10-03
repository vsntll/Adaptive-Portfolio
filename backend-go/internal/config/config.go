package config

import (
	"ioutil"	
	"time"

	"gopkg.in/yaml.v2"
)

// Config holds all configuration settings
type Config struct {
	WebDriver WebDriverConfig `yaml:"webdriver"`
	LinkedIn  LinkedInConfig  `yaml:"linkedin"`
	Export    ExportConfig    `yaml:"export"`
	Delays    DelayConfig     `yaml:"delays"`
	UserAgents []string       `yaml:"user_agents"`
	Logging   LoggingConfig   `yaml:"logging"`
}

// WebDriverConfig contains Selenium WebDriver settings
type WebDriverConfig struct {
	ChromeDriverPath string `yaml:"chromedriver_path"`
	Headless         bool   `yaml:"headless"`
	WindowSize       string `yaml:"window_size"`
	Port             int    `yaml:"port"`
	Timeout          int    `yaml:"timeout_seconds"`
}

// LinkedInConfig contains LinkedIn-specific settings
type LinkedInConfig struct {
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
	BaseURL  string `yaml:"base_url"`
	LoginURL string `yaml:"login_url"`
}

// ExportConfig contains data export settings
type ExportConfig struct {
	OutputDir string `yaml:"output_dir"`
	Format    string `yaml:"format"` // csv, json, both
}

// DelayConfig contains anti-detection delay settings
type DelayConfig struct {
	MinDelay      time.Duration `yaml:"min_delay"`
	MaxDelay      time.Duration `yaml:"max_delay"`
	PageLoad      time.Duration `yaml:"page_load_delay"`
	ElementSearch time.Duration `yaml:"element_search_delay"`
	ScrollDelay   time.Duration `yaml:"scroll_delay"`
}

// LoggingConfig contains logging settings
type LoggingConfig struct {
	Level    string `yaml:"level"`
	FilePath string `yaml:"file_path"`
	Console  bool   `yaml:"console"`
}

// Load reads configuration from YAML file
func Load(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	// Set defaults if not specified
	setDefaults(&config)

	return &config, nil
}

// setDefaults sets default values for configuration
func setDefaults(config *Config) {
	if config.WebDriver.Port == 0 {
		config.WebDriver.Port = 9515
	}
	
	if config.WebDriver.Timeout == 0 {
		config.WebDriver.Timeout = 30
	}
	
	if config.WebDriver.WindowSize == "" {
		config.WebDriver.WindowSize = "1920,1080"
	}
	
	if config.LinkedIn.BaseURL == "" {
		config.LinkedIn.BaseURL = "https://www.linkedin.com"
	}
	
	if config.LinkedIn.LoginURL == "" {
		config.LinkedIn.LoginURL = "https://www.linkedin.com/login"
	}
	
	if config.Export.OutputDir == "" {
		config.Export.OutputDir = "data/output"
	}
	
	if config.Export.Format == "" {
		config.Export.Format = "csv"
	}
	
	if config.Delays.MinDelay == 0 {
		config.Delays.MinDelay = 2 * time.Second
	}
	
	if config.Delays.MaxDelay == 0 {
		config.Delays.MaxDelay = 5 * time.Second
	}
	
	if config.Delays.PageLoad == 0 {
		config.Delays.PageLoad = 3 * time.Second
	}
	
	if config.Delays.ElementSearch == 0 {
		config.Delays.ElementSearch = 10 * time.Second
	}
	
	if config.Delays.ScrollDelay == 0 {
		config.Delays.ScrollDelay = 500 * time.Millisecond
	}
	
	if config.Logging.Level == "" {
		config.Logging.Level = "info"
	}
	
	if config.Logging.FilePath == "" {
		config.Logging.FilePath = "logs/scraper.log"
	}
	
	if len(config.UserAgents) == 0 {
		config.UserAgents = getDefaultUserAgents()
	}
}

// getDefaultUserAgents returns a list of default user agents
func getDefaultUserAgents() []string {
	return []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/119.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:109.0) Gecko/20100101 Firefox/119.0",
	}
}