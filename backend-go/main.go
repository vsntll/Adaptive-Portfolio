package main

import (
    "encoding/json"
    "net/http"
    "log"
)

// Define Portfolio data struct (customize fields as needed)
type Portfolio struct {
    Name        string   `json:"name"`
    Title       string   `json:"title"`
    Skills      []string `json:"skills"`
    Projects    []string `json:"projects"`
    Description string   `json:"description"`
}

// Mock portfolio data (replace with real data or DB as needed)
var mockPortfolio = Portfolio{
    Name:        "Avie Vasantlal",
    Title:       "Embedded Systems Hardware Engineer",
    Skills:      []string{"Go", "React", "C", "JavaScript"},
    Projects:    []string{"Adaptive-Portfolio", "Other Project"},
    Description: "Dynamic portfolio site powered by Go backend and React frontend.",
}

// Handler to serve portfolio data as JSON
func portfolioHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(mockPortfolio)
}

func main() {
    http.HandleFunc("/api/portfolio", portfolioHandler)

    // Enable CORS for frontend requests during development
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        if r.Method == "OPTIONS" {
            w.Header().Set("Access-Control-Allow-Methods", "GET")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
            w.WriteHeader(http.StatusNoContent)
            return
        }
    })

    log.Println("Go backend running at http://localhost:8080 ...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
