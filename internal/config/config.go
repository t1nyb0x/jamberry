package config

import (
	"fmt"
	"os"
	"strings"
)

// Config はアプリケーションの設定を保持します
type Config struct {
	DiscordBotToken  string
	TrackTasteAPIURL string
	RedisURL         string
	LogLevel         string
}

// Load は環境変数から設定を読み込みます
func Load() (*Config, error) {
	cfg := &Config{
		DiscordBotToken:  os.Getenv("DISCORD_BOT_TOKEN"),
		TrackTasteAPIURL: os.Getenv("TRACKTASTE_API_URL"),
		RedisURL:         os.Getenv("REDIS_URL"),
		LogLevel:         os.Getenv("LOG_LEVEL"),
	}

	// 必須項目のバリデーション
	var missing []string
	if cfg.DiscordBotToken == "" {
		missing = append(missing, "DISCORD_BOT_TOKEN")
	}
	if cfg.TrackTasteAPIURL == "" {
		missing = append(missing, "TRACKTASTE_API_URL")
	}
	if cfg.RedisURL == "" {
		missing = append(missing, "REDIS_URL")
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}

	// デフォルト値の設定
	if cfg.LogLevel == "" {
		cfg.LogLevel = "INFO"
	}

	// TrackTasteAPIURLの末尾スラッシュを除去
	cfg.TrackTasteAPIURL = strings.TrimSuffix(cfg.TrackTasteAPIURL, "/")

	return cfg, nil
}
