package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		os.Exit(1)
	}

	// Get the session cookie from the .env file
	sessionCookie := os.Getenv("SESSION")
	if sessionCookie == "" {
		fmt.Println("SESSION environment variable not set")
		os.Exit(1)
	}

	// Set the day of the Advent of Code
	day := 1

	// Construct the URL
	url := fmt.Sprintf("https://adventofcode.com/2024/day/%d/input", day)

	// Fetch the input data
	data, err := fetchData(url, sessionCookie)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		os.Exit(1)
	}

	// Parse the input into two slices
	left, right, err := parseInput(data)
	if err != nil {
		fmt.Println("Error parsing input:", err)
		os.Exit(1)
	}

	// Compute and print the similarity score
	similarityScore := computeSimilarityScore(left, right)
	fmt.Println("Similarity score:", similarityScore)
}

// fetchData fetches input data from the Advent of Code website
func fetchData(url, sessionCookie string) ([]byte, error) {
	// Create the HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add the session cookie to the request
	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: sessionCookie,
	})

	// Create the HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check if the HTTP response status is OK
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP request failed with status code %d, response: %s", resp.StatusCode, string(body))
	}

	// Read and return the response body
	return ioutil.ReadAll(resp.Body)
}

// parseInput parses the input data into two integer slices
func parseInput(data []byte) ([]int, []int, error) {
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	var left, right []int

	for _, line := range lines {
		numbers := strings.Fields(line)
		if len(numbers) != 2 {
			return nil, nil, fmt.Errorf("invalid input line: %s", line)
		}
		var l, r int
		_, err := fmt.Sscanf(numbers[0], "%d", &l)
		if err != nil {
			return nil, nil, err
		}
		_, err = fmt.Sscanf(numbers[1], "%d", &r)
		if err != nil {
			return nil, nil, err
		}
		left = append(left, l)
		right = append(right, r)
	}

	return left, right, nil
}

// computeSimilarityScore calculates the similarity score (Part Two)
func computeSimilarityScore(left, right []int) int {
	// Count occurrences of numbers in the right list
	rightCount := make(map[int]int)
	for _, num := range right {
		rightCount[num]++
	}

	// Calculate similarity score
	similarityScore := 0
	for _, num := range left {
		count := rightCount[num]
		similarityScore += num * count
	}

	return similarityScore
}
