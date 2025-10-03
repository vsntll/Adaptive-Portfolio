#!/bin/bash

# LinkedIn Scraper Setup Script
# This script sets up the development environment for the LinkedIn scraper

set -e  # Exit on any error

echo "🚀 Setting up LinkedIn Profile Scraper..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Go is installed
check_go() {
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed. Please install Go 1.19 or later."
        print_status "Visit: https://golang.org/doc/install"
        exit 1
    fi
    
    GO_VERSION=$(go version | cut -d' ' -f3 | sed 's/go//')
    print_success "Go $GO_VERSION is installed"
}

# Check if Chrome/Chromium is installed
check_chrome() {
    if command -v google-chrome &> /dev/null; then
        CHROME_VERSION=$(google-chrome --version | cut -d' ' -f3)
        print_success "Google Chrome $CHROME_VERSION is installed"
        return 0
    elif command -v chromium &> /dev/null; then
        CHROMIUM_VERSION=$(chromium --version | cut -d' ' -f2)
        print_success "Chromium $CHROMIUM_VERSION is installed"
        return 0
    elif command -v chromium-browser &> /dev/null; then
        CHROMIUM_VERSION=$(chromium-browser --version | cut -d' ' -f2)
        print_success "Chromium Browser $CHROMIUM_VERSION is installed"
        return 0
    else
        print_warning "Chrome/Chromium not found. Please install Chrome or Chromium browser."
        return 1
    fi
}

# Install ChromeDriver
install_chromedriver() {
    print_status "Checking for ChromeDriver..."
    
    if command -v chromedriver &> /dev/null; then
        CHROMEDRIVER_VERSION=$(chromedriver --version | cut -d' ' -f2)
        print_success "ChromeDriver $CHROMEDRIVER_VERSION is already installed"
        return 0
    fi
    
    print_status "ChromeDriver not found. Installing..."
    
    # Detect OS
    OS=$(uname -s)
    case $OS in
        "Darwin")
            # macOS
            if command -v brew &> /dev/null; then
                print_status "Installing ChromeDriver via Homebrew..."
                brew install chromedriver
            else
                print_error "Homebrew not found. Please install ChromeDriver manually:"
                print_status "1. Download from https://chromedriver.chromium.org/"
                print_status "2. Extract and place in /usr/local/bin/"
                print_status "3. Make executable: chmod +x /usr/local/bin/chromedriver"
                return 1
            fi
            ;;
        "Linux")
            # Linux
            print_status "Installing ChromeDriver for Linux..."
            
            # Get latest ChromeDriver version
            LATEST_VERSION=$(curl -s "https://chromedriver.storage.googleapis.com/LATEST_RELEASE")
            
            # Download ChromeDriver
            wget -q "https://chromedriver.storage.googleapis.com/$LATEST_VERSION/chromedriver_linux64.zip"
            unzip -q chromedriver_linux64.zip
            sudo mv chromedriver /usr/local/bin/
            sudo chmod +x /usr/local/bin/chromedriver
            rm chromedriver_linux64.zip
            
            print_success "ChromeDriver installed to /usr/local/bin/"
            ;;
        *)
            print_error "Unsupported operating system: $OS"
            print_status "Please install ChromeDriver manually from https://chromedriver.chromium.org/"
            return 1
            ;;
    esac
}

# Create directory structure
create_directories() {
    print_status "Creating directory structure..."
    
    directories=(
        "data/output"
        "logs"
        "bin"
    )
    
    for dir in "${directories[@]}"; do
        mkdir -p "$dir"
        print_success "Created directory: $dir"
    done
}

# Initialize Go module
init_go_module() {
    print_status "Initializing Go module..."
    
    if [ ! -f "go.mod" ]; then
        go mod init linkedin-scraper
        print_success "Go module initialized"
    else
        print_success "Go module already exists"
    fi
    
    print_status "Installing dependencies..."
    go get github.com/tebeka/selenium
    go get gopkg.in/yaml.v2
    go get github.com/gocarina/gocsv
    
    go mod tidy
    print_success "Dependencies installed"
}

# Create config file
create_config() {
    print_status "Creating configuration file..."
    
    if [ ! -f "configs/config.yaml" ]; then
        cp "configs/config.yaml" "configs/config.yaml.example" 2>/dev/null || true
        print_warning "Please update configs/config.yaml with your LinkedIn credentials"
    else
        print_success "Configuration file already exists"
    fi
}

# Build the application
build_app() {
    print_status "Building application..."
    
    go build -o bin/linkedin-scraper cmd/scraper/main.go
    
    if [ $? -eq 0 ]; then
        print_success "Application built successfully: bin/linkedin-scraper"
    else
        print_error "Build failed"
        exit 1
    fi
}

# Main setup function
main() {
    echo "========================================"
    echo "  LinkedIn Profile Scraper Setup"
    echo "========================================"
    echo ""
    
    check_go
    check_chrome
    install_chromedriver
    create_directories
    init_go_module
    create_config
    build_app
    
    echo ""
    echo "========================================"
    print_success "Setup completed successfully!"
    echo "========================================"
    echo ""
    
    print_status "Next steps:"
    echo "1. Update configs/config.yaml with your LinkedIn credentials"
    echo "2. Run the scraper:"
    echo "   ./bin/linkedin-scraper -profile 'https://www.linkedin.com/in/username/'"
    echo ""
    
    print_warning "Legal Notice:"
    echo "This tool is for educational purposes only."
    echo "Ensure compliance with LinkedIn's Terms of Service and applicable laws."
    echo "Use responsibly and respect rate limits."
    echo ""
}

# Run main function
main "$@"