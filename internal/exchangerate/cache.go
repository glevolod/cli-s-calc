package exchangerate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type CacheData struct {
	Base      string             `json:"base"`
	Rates     map[string]float64 `json:"rates"`
	Timestamp time.Time          `json:"timestamp"`
	Source    string             `json:"source"`
	ExpiresAt time.Time          `json:"expires_at"`
}

type Cache struct {
	cacheDir string
	ttl      time.Duration
}

func NewCache() (*Cache, error) {
	cacheDir, err := getCacheDir()
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	ttlHours := 24
	if ttlEnv := os.Getenv("S_CALC_CACHE_TTL"); ttlEnv != "" {
		if parsed, err := time.ParseDuration(ttlEnv + "h"); err == nil {
			ttlHours = int(parsed.Hours())
		}
	}

	return &Cache{
		cacheDir: cacheDir,
		ttl:      time.Duration(ttlHours) * time.Hour,
	}, nil
}

func getCacheDir() (string, error) {
	if customDir := os.Getenv("S_CALC_CACHE_DIR"); customDir != "" {
		return customDir, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	// Use .cache on Unix, LocalAppData on Windows
	if os.Getenv("LOCALAPPDATA") != "" {
		return filepath.Join(os.Getenv("LOCALAPPDATA"), "s-calc"), nil
	}
	return filepath.Join(homeDir, ".cache", "s-calc"), nil
}

func (c *Cache) Get(baseCurrency string) (*CacheData, error) {
	cacheFile := filepath.Join(c.cacheDir, fmt.Sprintf("rates-%s.json", baseCurrency))

	data, err := os.ReadFile(cacheFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil 
		}
		return nil, fmt.Errorf("failed to read cache file: %w", err)
	}

	var cacheData CacheData
	if err := json.Unmarshal(data, &cacheData); err != nil {
		return nil, fmt.Errorf("failed to parse cache file: %w", err)
	}

	if time.Now().After(cacheData.ExpiresAt) {
		return nil, nil
	}

	return &cacheData, nil
}

func (c *Cache) Set(baseCurrency string, rates map[string]float64, source string) error {
	cacheFile := filepath.Join(c.cacheDir, fmt.Sprintf("rates-%s.json", baseCurrency))

	now := time.Now()
	cacheData := CacheData{
		Base:      baseCurrency,
		Rates:     rates,
		Timestamp: now,
		Source:    source,
		ExpiresAt: now.Add(c.ttl),
	}

	data, err := json.MarshalIndent(cacheData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cache data: %w", err)
	}

	if err := os.WriteFile(cacheFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	return nil
}
