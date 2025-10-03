package util

import (
    "log"
)

// CheckError logs a fatal error if err is not nil
func CheckError(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

// ContainsString checks if a string slice contains a given string
func ContainsString(slice []string, str string) bool {
    for _, v := range slice {
        if v == str {
            return true
        }
    }
    return false
}

// RemoveDuplicates removes duplicate strings from a slice
func RemoveDuplicates(slice []string) []string {
    unique := make(map[string]bool)
    var result []string
    for _, v := range slice {
        if !unique[v] {
            unique[v] = true
            result = append(result, v)
        }
    }
    return result
}
