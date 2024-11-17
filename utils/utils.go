package utils

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
