package tracktaste

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_FetchHealth(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse func(w http.ResponseWriter, r *http.Request)
		wantErr        bool
		wantVersion    string
		wantStatus     string
		wantServices   int
	}{
		{
			name: "successful health check",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/healthz" {
					t.Errorf("unexpected path: %s", r.URL.Path)
				}
				resp := Response[HealthResult]{
					Status: 200,
					Result: HealthResult{
						Status:    "healthy",
						Version:   "1.0.0",
						BuildTime: "2024-01-01T00:00:00Z",
						GitCommit: "abc1234",
						Uptime:    "1h30m",
						Services: map[string]string{
							"spotify":  "enabled",
							"kkbox":    "enabled",
							"deezer":   "enabled",
							"redis":    "enabled",
							"lastfm":   "disabled",
						},
					},
				}
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(resp)
			},
			wantErr:      false,
			wantVersion:  "1.0.0",
			wantStatus:   "healthy",
			wantServices: 5,
		},
		{
			name: "empty services",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				resp := Response[HealthResult]{
					Status: 200,
					Result: HealthResult{
						Status:   "healthy",
						Version:  "1.0.0",
						Services: map[string]string{},
					},
				}
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(resp)
			},
			wantErr:      false,
			wantVersion:  "1.0.0",
			wantStatus:   "healthy",
			wantServices: 0,
		},
		{
			name: "server error",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				resp := APIError{
					Status:  500,
					Message: "Internal Server Error",
					Code:    "INTERNAL_ERROR",
				}
				_ = json.NewEncoder(w).Encode(resp)
			},
			wantErr: true,
		},
		{
			name: "service unavailable",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusServiceUnavailable)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.serverResponse))
			defer server.Close()

			client := NewClient(server.URL)
			health, err := client.FetchHealth(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("FetchHealth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if health.Version != tt.wantVersion {
					t.Errorf("FetchHealth() Version = %v, want %v", health.Version, tt.wantVersion)
				}
				if health.Status != tt.wantStatus {
					t.Errorf("FetchHealth() Status = %v, want %v", health.Status, tt.wantStatus)
				}
				if len(health.Services) != tt.wantServices {
					t.Errorf("FetchHealth() Services count = %v, want %v", len(health.Services), tt.wantServices)
				}
			}
		})
	}
}

func TestClient_FetchHealth_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 応答を遅延させる（キャンセルをテスト）
		<-r.Context().Done()
	}))
	defer server.Close()

	client := NewClient(server.URL)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // 即座にキャンセル

	_, err := client.FetchHealth(ctx)
	if err == nil {
		t.Error("FetchHealth() expected error for cancelled context, got nil")
	}
}
