package wow

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/Diony-source/CroupierBot/internal/config"
)

// --- Structs for Blizzard Auth ---
type BnetTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// --- NEW: Structs for Character Equipment ---
type Item struct {
	Level struct {
		Value float64 `json:"value"`
	} `json:"level"`
	Name string `json:"name"`
}
type EquippedItems struct {
	Item         Item `json:"item"`
	EquippedItemLevel float64 `json:"ilevel"`
}
type CharacterEquipmentSummary struct {
	EquippedItems     []EquippedItems `json:"equipped_items"`
	AverageItemLevel  float64         `json:"average_item_level"`
	EquippedItemLevel float64         `json:"equipped_item_level"`
	Character         struct {
		Name  string `json:"name"`
		Realm struct {
			Slug string `json:"slug"`
		} `json:"realm"`
	} `json:"character"`
	ActiveSpec struct {
		Name string `json:"name"`
	} `json:"active_spec"`
	CharacterClass struct {
		Name string `json:"name"`
	} `json:"character_class"`
}


// --- Blizzard API Client (No changes here) ---
type bnetClient struct {
	accessToken string
	tokenExpiry time.Time
	mutex       sync.Mutex
}
var client = &bnetClient{}

func (c *bnetClient) getAccessToken() (string, error) {
	// ... (Bu fonksiyonun içi aynı kalıyor)
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
	if err != nil { return "", fmt.Errorf("could not create token request: %w", err) }
	req.SetBasicAuth(config.Cfg.BlizzardClientID, config.Cfg.BlizzardClientSecret)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil { return "", fmt.Errorf("could not perform token request: %w", err) }
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK { return "", fmt.Errorf("blizzard auth API returned non-200 status: %d", resp.StatusCode) }
	var tokenResponse BnetTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil { return "", fmt.Errorf("could not decode token response: %w", err) }
	c.accessToken = tokenResponse.AccessToken
	c.tokenExpiry = time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second)
	fmt.Println("Successfully fetched new Blizzard access token.")
	return c.accessToken, nil
}

// --- NEW FUNCTION ---
// GetCharacterEquipment fetches a character's equipment summary from the Blizzard API.
func GetCharacterEquipment(characterName, serverName, region string) (*CharacterEquipmentSummary, error) {
	token, err := client.getAccessToken()
	if err != nil {
		return nil, fmt.Errorf("could not get blizzard access token: %w", err)
	}

	// Format server name for the URL (e.g., "Twisting Nether" -> "twisting-nether")
	serverSlug := strings.ToLower(strings.ReplaceAll(serverName, " ", "-"))
	charName := strings.ToLower(characterName)

	apiURL := fmt.Sprintf("https://%s.api.blizzard.com/profile/wow/character/%s/%s/equipment",
		region,
		serverSlug,
		charName,
	)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create equipment request: %w", err)
	}

	// Add necessary headers for the Blizzard API request
	q := req.URL.Query()
	q.Add("namespace", "profile-"+region)
	q.Add("locale", "en_US")
	req.URL.RawQuery = q.Encode()
	req.Header.Add("Authorization", "Bearer "+token)

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not perform equipment request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("blizzard equipment API returned non-200 status: %d", resp.StatusCode)
	}

	var summary CharacterEquipmentSummary
	if err := json.NewDecoder(resp.Body).Decode(&summary); err != nil {
		return nil, fmt.Errorf("could not decode equipment response: %w", err)
	}

	return &summary, nil
}