package helix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	ClientID    string
	AccessToken string
}

func NewHelixClient(clientID, accessToken string) *Client {
	return &Client{
		ClientID:    clientID,
		AccessToken: accessToken,
	}
}

type GetChannelResponse struct {
	Data []struct {
		BroadcasterId         string   `json:"broadcaster_id"`
		BroadcasterLogin      string   `json:"broadcaster_login"`
		BroadcasterName       string   `json:"broadcaster_name"`
		BroadcasterLanguage   string   `json:"broadcaster_language"`
		GameId                string   `json:"game_id"`
		GameName              string   `json:"game_name"`
		Title                 string   `json:"title"`
		Tags                  []string `json:"tags"`
		ContentClassification []string `json:"content_classification"`
	} `json:"data"`
}

func (c *Client) GetChannel(broadcasterID string) (*GetChannelResponse, error) {
	url := fmt.Sprintf("https://api.twitch.tv/helix/channels?broadcaster_id=%s", broadcasterID)

	token := fmt.Sprintf("Bearer %s", c.AccessToken)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Client-ID", c.ClientID)
	req.Header.Set("Authorization", token)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var channelResponse GetChannelResponse

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &channelResponse)
	if err != nil {
		return nil, err
	}

	return &channelResponse, nil
}

type UpdateChannelRequest struct {
	GameID                      string   `json:"game_id,omitempty"`
	BroadcasterLanguage         string   `json:"broadcaster_language,omitempty"`
	Title                       string   `json:"title,omitempty"`
	Delay                       string   `json:"delay,omitempty"`
	Tags                        []string `json:"tags,omitempty"`
	ContentClassificationLabels []struct {
		ID        string `json:"id"`
		IsEnabled bool   `json:"is_enabled"`
	} `json:"content_classification_labels,omitempty"`
}

func (c *Client) UpdateChannel(broadcasterId string, updateRequest UpdateChannelRequest) error {
	url := fmt.Sprintf("https://api.twitch.tv/helix/channels?broadcaster_id=%s", broadcasterId)

	token := fmt.Sprintf("Bearer %s", c.AccessToken)

	body, err := json.Marshal(updateRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Client-ID", c.ClientID)
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// return error and response body as string
	if resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d, %s", resp.StatusCode, string(body))
	}

	return nil
}

type GetGameResponse struct {
	Data []struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		BoxArtURL string `json:"box_art_url"`
		IGDBID    string `json:"igdb_id"`
	} `json:"data"`
}

func (c *Client) GetGameByName(gameName string) (*GetGameResponse, error) {
	gameName = strings.Join(strings.Split(gameName, " "), "+")
	url := fmt.Sprintf("https://api.twitch.tv/helix/games?name=%s", gameName)

	token := fmt.Sprintf("Bearer %s", c.AccessToken)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Client-ID", c.ClientID)
	req.Header.Set("Authorization", token)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d, %s", resp.StatusCode, body)
	}

	var gameResponse GetGameResponse

	err = json.Unmarshal(body, &gameResponse)
	if err != nil {
		return nil, err
	}

	return &gameResponse, nil
}

func (c *Client) GetGameById(gameId string) (*GetGameResponse, error) {
	url := fmt.Sprintf("https://api.twitch.tv/helix/games?id=%s", gameId)

	token := fmt.Sprintf("Bearer %s", c.AccessToken)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Client-ID", c.ClientID)
	req.Header.Set("Authorization", token)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var gameResponse GetGameResponse
	err = json.Unmarshal(body, &gameResponse)
	if err != nil {
		return nil, err
	}

	return &gameResponse, nil
}

type ChannelReward struct {
	BroadcasterName     string `json:"broadcaster_name"`
	BroadcasterLogin    string `json:"broadcaster_login"`
	BroadcasterID       string `json:"broadcaster_id"`
	ID                  string `json:"id"`
	Image               string `json:"image"`
	BackgroundColor     string `json:"background_color"`
	IsEnabled           bool   `json:"is_enabled"`
	Cost                int    `json:"cost"`
	Title               string `json:"title"`
	Prompt              string `json:"prompt"`
	IsUserInputRequired bool   `json:"is_user_input_required"`
	MaxPerStreamSetting struct {
		IsEnabled    bool `json:"is_enabled"`
		MaxPerStream int  `json:"max_per_stream"`
	} `json:"max_per_stream_setting"`
	MaxPerUserPerStreamSetting struct {
		IsEnabled           bool `json:"is_enabled"`
		MaxPerUserPerStream int  `json:"max_per_user_per_stream"`
	} `json:"max_per_user_per_stream_setting"`
	GlobalCooldownSetting struct {
		IsEnabled             bool `json:"is_enabled"`
		GlobalCooldownSeconds int  `json:"global_cooldown_seconds"`
	} `json:"global_cooldown_setting"`
	IsPaused     bool `json:"is_paused"`
	IsInStock    bool `json:"is_in_stock"`
	DefaultImage struct {
		URL1x string `json:"url_1x"`
		URL2x string `json:"url_2x"`
		URL4x string `json:"url_4x"`
	} `json:"default_image"`
	ShouldRedemptionsSkipRequestQueue bool   `json:"should_redemptions_skip_request_queue"`
	RedemptionsRedeemedCurrentStream  string `json:"redemptions_redeemed_current_stream"`
	CooldownExpiresAt                 string `json:"cooldown_expires_at"`
}

type GetChannelRewardsResponse struct {
	Data []ChannelReward `json:"data"`
}

func (c *Client) GetChannelRewards(broadcasterId string) (*[]ChannelReward, error) {
	url := fmt.Sprintf("https://api.twitch.tv/helix/channel_points/custom_rewards?broadcaster_id=%s", broadcasterId)

	token := fmt.Sprintf("Bearer %s", c.AccessToken)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Client-ID", c.ClientID)
	req.Header.Set("Authorization", token)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rewardsResponse GetChannelRewardsResponse
	err = json.Unmarshal(body, &rewardsResponse)
	if err != nil {
		return nil, err
	}

	return &rewardsResponse.Data, nil
}

func (c *Client) GetChannelRewardByName(broadcasterId, rewardName string) (*ChannelReward, error) {
	rewards, err := c.GetChannelRewards(broadcasterId)
	if err != nil {
		return nil, err
	}

	for _, reward := range *rewards {
		if reward.Title == rewardName {
			return &reward, nil
		}
	}

	return nil, nil
}

func (c *Client) GetChannelRewardByID(broadcasterId, rewardID string) (*ChannelReward, error) {
	rewards, err := c.GetChannelRewards(broadcasterId)
	if err != nil {
		return nil, err
	}

	for _, reward := range *rewards {
		if reward.ID == rewardID {
			return &reward, nil
		}
	}

	return nil, nil
}

type CreateChannelRewardRequest struct {
	Title                   string `json:"title"`
	Prompt                  string `json:"prompt,omitempty"`
	Cost                    int    `json:"cost"`
	BackgroundColor         string `json:"background_color,omitempty"`
	IsGlobalCooldownEnabled bool   `json:"is_global_cooldown_enabled,omitempty"`
	GlobalCooldownSeconds   int    `json:"global_cooldown_seconds,omitempty"`
	IsEnabled               bool   `json:"is_enabled"`
}

func (c *Client) CreateChannelReward(broadcasterId string, rewardRequest *CreateChannelRewardRequest) (*ChannelReward, error) {
	url := fmt.Sprintf("https://api.twitch.tv/helix/channel_points/custom_rewards?broadcaster_id=%s", broadcasterId)

	token := fmt.Sprintf("Bearer %s", c.AccessToken)

	body, err := json.Marshal(rewardRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Client-ID", c.ClientID)
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rewardResponse GetChannelRewardsResponse
	err = json.Unmarshal(body, &rewardResponse)
	if err != nil {
		return nil, err
	}

	return &rewardResponse.Data[0], nil
}

type UpdateChannelRewardRequest struct {
	Title                   string `json:"title"`
	Prompt                  string `json:"prompt,omitempty"`
	Cost                    int    `json:"cost"`
	BackgroundColor         string `json:"background_color,omitempty"`
	IsGlobalCooldownEnabled bool   `json:"is_global_cooldown_enabled,omitempty"`
	GlobalCooldownSeconds   int    `json:"global_cooldown_seconds,omitempty"`
	IsEnabled               bool   `json:"is_enabled"`
}

func (c *Client) UpdateChannelReward(broadcasterID, rewardID string, updateRequest *UpdateChannelRewardRequest) (*ChannelReward, error) {
	url := fmt.Sprintf("https://api.twitch.tv/helix/channel_points/custom_rewards?broadcaster_id=%s&id=%s", broadcasterID, rewardID)

	token := fmt.Sprintf("Bearer %s", c.AccessToken)

	body, err := json.Marshal(&updateRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Client-ID", c.ClientID)
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d, %s", resp.StatusCode, string(body))
	}

	var rewardResponse GetChannelRewardsResponse
	err = json.Unmarshal(body, &rewardResponse)
	if err != nil {
		return nil, err
	}

	return &rewardResponse.Data[0], nil
}

func (c *Client) DeleteChannelReward(broadcasterID, rewardID string) error {
	url := fmt.Sprintf("https://api.twitch.tv/helix/channel_points/custom_rewards?broadcaster_id=%s&id=%s", broadcasterID, rewardID)

	token := fmt.Sprintf("Bearer %s", c.AccessToken)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Client-ID", c.ClientID)
	req.Header.Set("Authorization", token)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
