// internal/wow/raiderio.go
package wow

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// AffixDetails struct holds the details of a single affix.
type AffixDetails struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// MythicPlusAffixesResponse struct matches the structure of the Raider.IO API response.
type MythicPlusAffixesResponse struct {
	Title        string         `json:"title"`
	AffixDetails []AffixDetails `json:"affix_details"`
}

// GetCurrentAffixes fetches the current Mythic+ affixes from the Raider.IO API.
func GetCurrentAffixes() (*MythicPlusAffixesResponse, error) {
	// Raider.IO API URL for current affixes. We specify our region, e.g., "eu".
	url := "https://raider.io/api/v1/mythic-plus/affixes?region=eu&locale=en"

	// Make the HTTP GET request.
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not make request to Raider.IO: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("raider.IO API returned a non-200 status code: %d", resp.StatusCode)
	}

	// Read the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %w", err)
	}

	// Unmarshal the JSON response into our struct.
	var affixes MythicPlusAffixesResponse
	if err := json.Unmarshal(body, &affixes); err != nil {
		return nil, fmt.Errorf("could not unmarshal JSON response: %w", err)
	}

	return &affixes, nil
}
