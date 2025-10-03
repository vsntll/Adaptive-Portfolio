package utils

import (
	"math/rand"
	"time"
)

// RandomSleep sleeps for a random duration between min and max
func RandomSleep(min, max time.Duration) {
	if max <= min {
		time.Sleep(min)
		return
	}
	
	diff := int64(max - min)
	randomDelay := time.Duration(rand.Int63n(diff)) + min
	time.Sleep(randomDelay)
}

// RandomInt generates a random integer between min and max (inclusive)
func RandomInt(min, max int) int {
	if max <= min {
		return min
	}
	return rand.Intn(max-min+1) + min
}

// RandomFloat generates a random float64 between min and max
func RandomFloat(min, max float64) float64 {
	if max <= min {
		return min
	}
	return rand.Float64()*(max-min) + min
}

// RandomBool returns a random boolean value
func RandomBool() bool {
	return rand.Intn(2) == 1
}

// RandomChoice returns a random element from the slice
func RandomChoice(choices []string) string {
	if len(choices) == 0 {
		return ""
	}
	return choices[rand.Intn(len(choices))]
}

// ShuffleStrings shuffles a slice of strings in place
func ShuffleStrings(slice []string) {
	rand.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
}

// GenerateRandomUserAgent generates a random user agent string
func GenerateRandomUserAgent() string {
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/119.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:109.0) Gecko/20100101 Firefox/119.0",
		"Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/119.0",
	}
	
	return RandomChoice(userAgents)
}

// init initializes the random seed
func init() {
	rand.Seed(time.Now().UnixNano())
}