package exchangerate

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ExchangeRateAPI struct {
	cache *Cache
}

func NewExchangeRateAPI() (*ExchangeRateAPI, error) {
	cache, err := NewCache()
	if err != nil {
		return nil, err
	}

	return &ExchangeRateAPI{
		cache: cache,
	}, nil
}

type RateResponse struct {
	Base  string             `json:"base"`
	Rates map[string]float64 `json:"rates"`
	Date  string             `json:"date"`
}

func (api *ExchangeRateAPI) GetRates(baseCurrency string) (map[string]float64, *RateInfo, error) {
	if cached, err := api.cache.Get(baseCurrency); err == nil && cached != nil {
		return cached.Rates, &RateInfo{
			Source:    cached.Source,
			Timestamp: cached.Timestamp,
			ExpiresAt: cached.ExpiresAt,
		}, nil
	}

	rates, info, err := api.fetchFromPrimary(baseCurrency)
	if err == nil {
		_ = api.cache.Set(baseCurrency, rates, info.Source)
		return rates, info, nil
	}

	rates, info, err = api.fetchFromFallback(baseCurrency)
	if err == nil {
		_ = api.cache.Set(baseCurrency, rates, info.Source)
		return rates, info, nil
	}

	if cached, err := api.cache.Get(baseCurrency); err == nil && cached != nil {
		return cached.Rates, &RateInfo{
			Source:    cached.Source + " (expired)",
			Timestamp: cached.Timestamp,
			ExpiresAt: cached.ExpiresAt,
		}, nil
	}

	return nil, nil, fmt.Errorf("failed to fetch exchange rates: %w", err)
}

type RateInfo struct {
	Source    string
	Timestamp time.Time
	ExpiresAt time.Time
}

func (api *ExchangeRateAPI) fetchFromPrimary(baseCurrency string) (map[string]float64, *RateInfo, error) {
	url := fmt.Sprintf("https://api.exchangerate-api.com/v4/latest/%s", baseCurrency)

	resp, err := http.Get(url)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	var rateResp RateResponse
	if err := json.Unmarshal(body, &rateResp); err != nil {
		return nil, nil, err
	}

	if rateResp.Rates == nil {
		rateResp.Rates = make(map[string]float64)
	}
	rateResp.Rates[baseCurrency] = 1.0

	return rateResp.Rates, &RateInfo{
		Source:    "exchangerate-api.com",
		Timestamp: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}, nil
}

func (api *ExchangeRateAPI) fetchFromFallback(baseCurrency string) (map[string]float64, *RateInfo, error) {
	url := fmt.Sprintf("https://api.exchangerate.host/latest?base=%s", baseCurrency)

	resp, err := http.Get(url)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	var response struct {
		Success bool               `json:"success"`
		Base    string             `json:"base"`
		Rates   map[string]float64 `json:"rates"`
		Date    string             `json:"date"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, nil, err
	}

	if !response.Success {
		return nil, nil, fmt.Errorf("API returned success=false")
	}

	if response.Rates == nil {
		response.Rates = make(map[string]float64)
	}
	response.Rates[baseCurrency] = 1.0

	return response.Rates, &RateInfo{
		Source:    "exchangerate.host",
		Timestamp: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}, nil
}
