package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/t1nyb0x/jamberry/internal/infrastructure/tracktaste"
)

// mockTracktasteClient はTracktasteClientのモック実装です
type mockTracktasteClient struct {
	fetchHealthFunc func(ctx context.Context) (*tracktaste.HealthResult, error)
}

func (m *mockTracktasteClient) FetchHealth(ctx context.Context) (*tracktaste.HealthResult, error) {
	if m.fetchHealthFunc != nil {
		return m.fetchHealthFunc(ctx)
	}
	return nil, errors.New("not implemented")
}

func TestMockTracktasteClient_FetchHealth(t *testing.T) {
	tests := []struct {
		name        string
		mockFunc    func(ctx context.Context) (*tracktaste.HealthResult, error)
		wantErr     bool
		wantVersion string
	}{
		{
			name: "successful health check",
			mockFunc: func(ctx context.Context) (*tracktaste.HealthResult, error) {
				return &tracktaste.HealthResult{
					Status:    "healthy",
					Version:   "1.0.0",
					BuildTime: "2024-01-01T00:00:00Z",
					GitCommit: "abc1234",
					Uptime:    "1h30m",
					Services: map[string]string{
						"spotify": "enabled",
						"kkbox":   "enabled",
					},
				}, nil
			},
			wantErr:     false,
			wantVersion: "1.0.0",
		},
		{
			name: "connection error",
			mockFunc: func(ctx context.Context) (*tracktaste.HealthResult, error) {
				return nil, errors.New("connection refused")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &mockTracktasteClient{
				fetchHealthFunc: tt.mockFunc,
			}

			result, err := client.FetchHealth(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("FetchHealth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if result.Version != tt.wantVersion {
					t.Errorf("FetchHealth() Version = %v, want %v", result.Version, tt.wantVersion)
				}
			}
		})
	}
}

func TestHealthResult_Services(t *testing.T) {
	result := &tracktaste.HealthResult{
		Status:  "healthy",
		Version: "1.0.0",
		Services: map[string]string{
			"spotify":       "enabled",
			"kkbox":         "enabled",
			"deezer":        "enabled",
			"musicbrainz":   "enabled",
			"lastfm":        "disabled",
			"youtube_music": "disabled",
			"redis":         "enabled",
		},
	}

	// サービス数の確認
	if len(result.Services) != 7 {
		t.Errorf("Expected 7 services, got %d", len(result.Services))
	}

	// enabled/disabled のカウント
	enabledCount := 0
	disabledCount := 0
	for _, status := range result.Services {
		if status == "enabled" {
			enabledCount++
		} else if status == "disabled" {
			disabledCount++
		}
	}

	if enabledCount != 5 {
		t.Errorf("Expected 5 enabled services, got %d", enabledCount)
	}
	if disabledCount != 2 {
		t.Errorf("Expected 2 disabled services, got %d", disabledCount)
	}
}
