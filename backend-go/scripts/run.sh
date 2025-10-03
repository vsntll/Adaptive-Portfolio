#!/bin/bash

# LinkedIn Scraper Run Script
# This script provides easy commands to run the LinkedIn scraper

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

# Default values
PROFILE_URL=""
OUTPUT_FORMAT="csv"
CONFIG_FILE="configs/config.yaml"
OUTPUT_FILE="data/output/profile.csv"
BINARY="bin/linkedin-scraper"

# Show usage information
show_usage() {
    echo "LinkedIn Profile Scraper - Run Script"
    echo ""
    echo "Usage: $0 [OPTIONS] PROFILE_URL"
    echo ""
    echo "Options:"
    echo "  -f, --format FORMAT    Output format: csv, json, or both (default: csv)"
    echo "  -o, --output FILE      Output file path (default: data/output/profile.csv)"
    echo "  -c, --config FILE      Configuration file path (default: configs/config.yaml)"
    echo "  -h, --help             Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 https://www.linkedin.com/in/username/"
    echo "  $0 -f json -o profiles.json https://www.linkedin.com/in/username/"
    echo "  $0 --format both --output data/profiles.csv https://www.linkedin.com/in/username/"
    echo ""
}

# Parse command line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -f|--format)
                OUTPUT_FORMAT="$2"
                shift 2
                ;;
            -o|--output)
                OUTPUT_FILE="$2"
                shift 2
                ;;
            -c|--config)
                CONFIG_FILE="$2"
                shift 2
                ;;
            -h|--help)
                show_usage
                exit 0
                ;;
            -*|--*)
                print_error "Unknown option $1"
                show_usage
                exit 1
                ;;
            *)
                if [ -z "$PROFILE_URL" ]; then
                    PROFILE_URL="$1"
                else
                    print_error "Multiple profile URLs provided"
                    show_usage
                    exit 1
                fi
                shift
                ;;
        esac
    done
}

# Validate inputs
validate_inputs() {
    if [ -z "$PROFILE_URL" ]; then
        print_error "Profile URL is required"
        show_usage
        exit 1
    fi
    
    # Validate profile URL format
    if [[ ! "$PROFILE_URL" =~ ^https?://.*linkedin\.com/in/[a-zA-Z0-9\-]+/?$ ]]; then
        print_error "Invalid LinkedIn profile URL format"
        print_status "Expected format: https://www.linkedin.com/in/username/"
        exit 1
    fi
    
    # Validate output format
    case $OUTPUT_FORMAT in
        csv|json|both)
            ;;
        *)
            print_error "Invalid output format: $OUTPUT_FORMAT"
            print_status "Supported formats: csv, json, both"
            exit 1
            ;;
    esac
    
    # Check if config file exists
    if [ ! -f "$CONFIG_FILE" ]; then
        print_error "Configuration file not found: $CONFIG_FILE"
        print_status "Please run 'make setup' or create the configuration file manually"
        exit 1
    fi
    
    # Check if binary exists
    if [ ! -f "$BINARY" ]; then
        print_error "Binary not found: $BINARY"
        print_status "Please run 'make build' to build the application"
        exit 1
    fi
}

# Create output directory
create_output_dir() {
    OUTPUT_DIR=$(dirname "$OUTPUT_FILE")
    if [ ! -d "$OUTPUT_DIR" ]; then
        print_status "Creating output directory: $OUTPUT_DIR"
        mkdir -p "$OUTPUT_DIR"
    fi
}

# Run the scraper
run_scraper() {
    print_status "Starting LinkedIn profile scraper..."
    print_status "Profile URL: $PROFILE_URL"
    print_status "Output format: $OUTPUT_FORMAT"
    print_status "Output file: $OUTPUT_FILE"
    print_status "Configuration: $CONFIG_FILE"
    echo ""
    
    # Build the command
    CMD="$BINARY -profile \"$PROFILE_URL\" -format \"$OUTPUT_FORMAT\" -output \"$OUTPUT_FILE\" -config \"$CONFIG_FILE\""
    
    print_status "Executing: $CMD"
    echo ""
    
    # Run the scraper
    if eval "$CMD"; then
        echo ""
        print_success "Scraping completed successfully!"
        
        # Show output files
        case $OUTPUT_FORMAT in
            csv)
                if [ -f "$OUTPUT_FILE" ]; then
                    print_success "CSV file created: $OUTPUT_FILE"
                    print_status "File size: $(du -h "$OUTPUT_FILE" | cut -f1)"
                fi
                ;;
            json)
                JSON_FILE="${OUTPUT_FILE%.*}.json"
                if [ -f "$JSON_FILE" ]; then
                    print_success "JSON file created: $JSON_FILE"
                    print_status "File size: $(du -h "$JSON_FILE" | cut -f1)"
                fi
                ;;
            both)
                if [ -f "$OUTPUT_FILE" ]; then
                    print_success "CSV file created: $OUTPUT_FILE"
                    print_status "File size: $(du -h "$OUTPUT_FILE" | cut -f1)"
                fi
                JSON_FILE="${OUTPUT_FILE%.*}.json"
                if [ -f "$JSON_FILE" ]; then
                    print_success "JSON file created: $JSON_FILE"
                    print_status "File size: $(du -h "$JSON_FILE" | cut -f1)"
                fi
                ;;
        esac
    else
        print_error "Scraping failed with exit code $?"
        exit 1
    fi
}

# Show profile preview
show_preview() {
    if [ -f "$OUTPUT_FILE" ] && [ "$OUTPUT_FORMAT" != "json" ]; then
        echo ""
        print_status "Profile preview:"
        echo "----------------------------------------"
        
        # Show first few lines of CSV (skip header)
        if command -v column &> /dev/null; then
            head -n 2 "$OUTPUT_FILE" | column -t -s ','
        else
            head -n 2 "$OUTPUT_FILE"
        fi
        
        echo "----------------------------------------"
    fi
}

# Main function
main() {
    echo "========================================"
    echo "  LinkedIn Profile Scraper"
    echo "========================================"
    echo ""
    
    parse_args "$@"
    validate_inputs
    create_output_dir
    run_scraper
    show_preview
    
    echo ""
    print_warning "Legal Notice:"
    print_warning "This tool is for educational purposes only."
    print_warning "Ensure compliance with LinkedIn's Terms of Service."
    echo ""
}

# Run main function with all arguments
main "$@"