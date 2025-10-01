package wow

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// --- Affix Structs (No changes) ---
type AffixDetails struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
type MythicPlusAffixesResponse struct {
	Title        string         `json:"title"`
	AffixDetails []AffixDetails `json:"affix_details"`
}

// --- Rio Structs (UPDATED) ---
type RaidProgress struct {
	Summary string `json:"summary"`
}
type MythicPlusScores struct {
	All float64 `json:"all"`
}

// This new struct correctly represents one season's score data.
type MythicPlusScoresBySeason struct {
	Season string           `json:"season"`
	Scores MythicPlusScores `json:"scores"`
}
type CharacterProfileResponse struct {
	Name            string                  `json:"name"`
	Race            string                  `json:"race"`
	Class           string                  `json:"class"`
	ActiveSpecName  string                  `json:"active_spec_name"`
	ThumbnailURL    string                  `json:"thumbnail_url"`
	ProfileURL      string                  `json:"profile_url"`
	RaidProgression map[string]RaidProgress `json:"raid_progression"`
	// This field is now a slice (array) to match the API response.
	MythicPlusScoresBySeason []MythicPlusScoresBySeason `json:"mythic_plus_scores_by_season"`
}

// --- Functions (No changes in GetCurrentAffixes) ---
func GetCurrentAffixes() (*MythicPlusAffixesResponse, error) {
	// ... (Bu fonksiyonun içi aynı kalıyor)
	url := "https://raider.io/api/v1/mythic-plus/affixes?region=eu&locale=en"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not make request to Raider.IO: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("raider.IO API returned a non-200 status code: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %w", err)
	}
	var affixes MythicPlusAffixesResponse
	if err := json.Unmarshal(body, &affixes); err != nil {
		return nil, fmt.Errorf("could not unmarshal JSON response: %w", err)
	}
	return &affixes, nil
}

// --- GetCharacterProfile (No changes) ---
func GetCharacterProfile(characterName, serverName, region string) (*CharacterProfileResponse, error) {
	// ... (Bu fonksiyonun içi aynı kalıyor)
	apiURL := fmt.Sprintf("https://raider.io/api/v1/characters/profile?region=%s&realm=%s&name=%s&fields=raid_progression,mythic_plus_scores_by_season:current",
		url.QueryEscape(region),
		url.QueryEscape(serverName),
		url.QueryEscape(characterName),
	)
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("could not make request to Raider.IO: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("raider.IO API returned a non-200 status code: %d. Character might not exist or server name is wrong", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %w", err)
	}
	var profile CharacterProfileResponse
	if err := json.Unmarshal(body, &profile); err != nil {
		return nil, fmt.Errorf("could not unmarshal JSON response: %w", err)
	}
	return &profile, nil
}
