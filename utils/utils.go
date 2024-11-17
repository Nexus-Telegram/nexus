package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
)

type QueryID struct {
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	// Include other fields if necessary
}

func MergeHeaders(commonHeaders map[string]string, queryIDheaders map[string]string) map[string]string {
	// Create a copy of the commonHeaders to avoid modifying the original map
	mergedHeaders := make(map[string]string)

	// Copy commonHeaders into mergedHeaders
	for key, value := range commonHeaders {
		mergedHeaders[key] = value
	}

	// Add or overwrite with queryIDheaders
	for key, value := range queryIDheaders {
		mergedHeaders[key] = value
	}

	return mergedHeaders
}

func parseQueryIDs() []string {
	file, err := os.Open("query_ids.txt")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)

	// Create a slice to store all the queryIDs
	var queryIDs []string

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" { // Ignore empty lines
			queryIDs = append(queryIDs, line)
		}
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	return queryIDs
}

func extractFullName(queryString string) (string, error) {
	// Parse the query string
	parsedQuery, err := url.ParseQuery(queryString)
	if err != nil {
		return "", fmt.Errorf("failed to parse query string: %w", err)
	}

	// Extract the 'user' parameter (it's URL-encoded JSON)
	encodedUser := parsedQuery.Get("user")
	if encodedUser == "" {
		return "", fmt.Errorf("user parameter not found in query string")
	}

	// Decode the URL-encoded string to get the raw JSON
	decodedUser, err := url.QueryUnescape(encodedUser)
	if err != nil {
		return "", fmt.Errorf("failed to decode user data: %w", err)
	}

	// Unmarshal the JSON into the User struct
	var user QueryID
	err = json.Unmarshal([]byte(decodedUser), &user)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal user data: %w", err)
	}

	// Use the username if it's available, otherwise use the first and last name
	if user.Username != "" {
		return user.Username, nil
	}

	// Combine first name and last name into a full name
	fullName := fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	// If last name is empty, return only the first name
	if user.LastName == "" {
		fullName = user.FirstName
	}

	// Return the full name
	return fullName, nil
}
