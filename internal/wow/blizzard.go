package wow

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/Diony-source/CroupierBot/internal/config"
)

// --- Structs for Blizzard Auth (No changes) ---
type BnetTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// --- Structs for Character Summary (No changes) ---
type CharacterSummary struct {
	Name              string  `json:"name"`
	AverageItemLevel  float64 `json:"average_item_level"`
	EquippedItemLevel float64 `json:"equipped_item_level"`
	ActiveSpec        struct {
		Name string `json:"name"`
	} `json:"active_spec"`
	CharacterClass struct {
		Name string `json:"name"`
	} `json:"character_class"`
	Realm struct {
		Slug string `json:"slug"`
	} `json:"realm"`
}

// --- Structs for Character Media (No changes) ---
type Asset struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
type CharacterMediaSummary struct {
	Assets []Asset `json:"assets"`
}

// --- STRUCTS FOR CHARACTER EQUIPMENT (CORRECTED) ---
type EquippedItem struct {
	Slot struct {
		Name string `json:"name"`
	} `json:"slot"`
	Name  string `json:"name"`
	Level struct {
		Value float64 `json:"value"`
	} `json:"level"`
	Enchantments []struct {
		DisplayString string `json:"display_string"`
	} `json:"enchantments"`
}
type CharacterEquipment struct {
	EquippedItems []EquippedItem `json:"equipped_items"`
}

// --- Blizzard API Client and Functions (No logical changes, just ensuring they are complete) ---
type bnetClient struct {
	accessToken string
	tokenExpiry time.Time
	mutex       sync.Mutex
}

var client = &bnetClient{}

func (c *bnetClient) getAccessToken() (string, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.accessToken != "" && time.Now().Before(c.tokenExpiry.Add(-1*time.Minute)) {
		return c.accessToken, nil
	}
	fmt.Println("Blizzard access token is expired or missing. Fetching a new one...")
	authURL := "https://oauth.battle.net/token"
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	req, err := http.NewRequest("POST", authURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("could not create token request: %w", err)
	}
	req.SetBasicAuth(config.Cfg.BlizzardClientID, config.Cfg.BlizzardClientSecret)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("could not perform token request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("blizzard auth API returned non-200 status: %d", resp.StatusCode)
	}
	var tokenResponse BnetTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", fmt.Errorf("could not decode token response: %w", err)
	}
	c.accessToken = tokenResponse.AccessToken
	c.tokenExpiry = time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second)
	fmt.Println("Successfully fetched new Blizzard access token.")
	return c.accessToken, nil
}

func makeBnetAPIRequest(apiURL, region string, target interface{}) error {
	token, err := client.getAccessToken()
	if err != nil {
		return err
	}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}
	q := req.URL.Query()
	q.Add("namespace", "profile-"+region)
	q.Add("locale", "en_US")
	req.URL.RawQuery = q.Encode()
	req.Header.Add("Authorization", "Bearer "+token)
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("could not perform request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned non-200 status: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read response body: %w", err)
	}
	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("could not decode response: %w", err)
	}
	return nil
}

func GetCharacterSummary(characterName, serverName, region string) (*CharacterSummary, error) {
	serverSlug := strings.ToLower(strings.ReplaceAll(serverName, " ", "-"))
	charName := strings.ToLower(characterName)
	apiURL := fmt.Sprintf("https://%s.api.blizzard.com/profile/wow/character/%s/%s", region, serverSlug, charName)
	var summary CharacterSummary
	err := makeBnetAPIRequest(apiURL, region, &summary)
	if err != nil {
		return nil, err
	}
	return &summary, nil
}

func GetCharacterMedia(characterName, serverName, region string) (*CharacterMediaSummary, error) {
	serverSlug := strings.ToLower(strings.ReplaceAll(serverName, " ", "-"))
	charName := strings.ToLower(characterName)
	apiURL := fmt.Sprintf("https://%s.api.blizzard.com/profile/wow/character/%s/%s/character-media", region, serverSlug, charName)
	var media CharacterMediaSummary
	err := makeBnetAPIRequest(apiURL, region, &media)
	if err != nil {
		return nil, err
	}
	return &media, nil
}

func GetCharacterEquipment(characterName, serverName, region string) (*CharacterEquipment, error) {
	serverSlug := strings.ToLower(strings.ReplaceAll(serverName, " ", "-"))
	charName := strings.ToLower(characterName)
	apiURL := fmt.Sprintf("https://%s.api.blizzard.com/profile/wow/character/%s/%s/equipment", region, serverSlug, charName)
	var equipment CharacterEquipment
	err := makeBnetAPIRequest(apiURL, region, &equipment)
	if err != nil {
		return nil, err
	}
	return &equipment, nil
}
