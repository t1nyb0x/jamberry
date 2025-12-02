package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/t1nyb0x/jamberry/internal/domain"
)

func TestRecommendUseCase_GetRecommend(t *testing.T) {
	tests := []struct {
		name            string
		input           RecommendInput
		fetchRecommend  func(ctx context.Context, spotifyURL string, mode domain.RecommendMode, limit int) (*domain.RecommendResult, error)
		fetchTrack      func(ctx context.Context, spotifyURL string) (*domain.Track, error)
		wantErr         bool
		errType         string
		wantMode        domain.RecommendMode
		wantItemsCount  int
	}{
		{
			name: "valid URL with default mode",
			input: RecommendInput{
				Input: "https://open.spotify.com/track/4iV5W9uYEdYUVa79Axb7Rh",
			},
			fetchRecommend: func(ctx context.Context, spotifyURL string, mode domain.RecommendMode, limit int) (*domain.RecommendResult, error) {
				return &domain.RecommendResult{
					SeedTrack: domain.Track{ID: "seed", Name: "Seed Track"},
					Items: []domain.SimilarTrack{
						{ID: "1", Name: "Similar 1"},
						{ID: "2", Name: "Similar 2"},
					},
					Mode: domain.RecommendModeBalanced,
				}, nil
			},
			fetchTrack: func(ctx context.Context, spotifyURL string) (*domain.Track, error) {
				return &domain.Track{ID: "seed", Name: "Seed Track", URL: spotifyURL}, nil
			},
			wantErr:        false,
			wantMode:       domain.RecommendModeBalanced,
			wantItemsCount: 2,
		},
		{
			name: "valid URL with similar mode",
			input: RecommendInput{
				Input: "https://open.spotify.com/track/4iV5W9uYEdYUVa79Axb7Rh",
				Mode:  domain.RecommendModeSimilar,
			},
			fetchRecommend: func(ctx context.Context, spotifyURL string, mode domain.RecommendMode, limit int) (*domain.RecommendResult, error) {
				if mode != domain.RecommendModeSimilar {
					return nil, errors.New("wrong mode")
				}
				return &domain.RecommendResult{
					SeedTrack: domain.Track{ID: "seed", Name: "Seed Track"},
					Items: []domain.SimilarTrack{
						{ID: "1", Name: "Similar 1"},
					},
					Mode: domain.RecommendModeSimilar,
				}, nil
			},
			fetchTrack: func(ctx context.Context, spotifyURL string) (*domain.Track, error) {
				return &domain.Track{ID: "seed", Name: "Seed Track"}, nil
			},
			wantErr:        false,
			wantMode:       domain.RecommendModeSimilar,
			wantItemsCount: 1,
		},
		{
			name: "valid URL with related mode",
			input: RecommendInput{
				Input: "https://open.spotify.com/track/4iV5W9uYEdYUVa79Axb7Rh",
				Mode:  domain.RecommendModeRelated,
			},
			fetchRecommend: func(ctx context.Context, spotifyURL string, mode domain.RecommendMode, limit int) (*domain.RecommendResult, error) {
				return &domain.RecommendResult{
					SeedTrack: domain.Track{ID: "seed", Name: "Seed Track"},
					Items: []domain.SimilarTrack{
						{ID: "1", Name: "Related 1"},
					},
					Mode: domain.RecommendModeRelated,
				}, nil
			},
			fetchTrack: func(ctx context.Context, spotifyURL string) (*domain.Track, error) {
				return &domain.Track{ID: "seed", Name: "Seed Track"}, nil
			},
			wantErr:        false,
			wantMode:       domain.RecommendModeRelated,
			wantItemsCount: 1,
		},
		{
			name: "valid URI",
			input: RecommendInput{
				Input: "spotify:track:4iV5W9uYEdYUVa79Axb7Rh",
			},
			fetchRecommend: func(ctx context.Context, spotifyURL string, mode domain.RecommendMode, limit int) (*domain.RecommendResult, error) {
				return &domain.RecommendResult{
					SeedTrack: domain.Track{ID: "seed", Name: "Seed Track"},
					Items: []domain.SimilarTrack{
						{ID: "1", Name: "Similar 1"},
					},
					Mode: domain.RecommendModeBalanced,
				}, nil
			},
			fetchTrack: func(ctx context.Context, spotifyURL string) (*domain.Track, error) {
				return &domain.Track{ID: "seed", Name: "Seed Track"}, nil
			},
			wantErr:        false,
			wantItemsCount: 1,
		},
		{
			name: "empty input",
			input: RecommendInput{
				Input: "",
			},
			wantErr: true,
			errType: "validation",
		},
		{
			name: "invalid URL - not spotify",
			input: RecommendInput{
				Input: "https://example.com/track/xxx",
			},
			wantErr: true,
			errType: "validation",
		},
		{
			name: "artist URL instead of track",
			input: RecommendInput{
				Input: "https://open.spotify.com/artist/4iV5W9uYEdYUVa79Axb7Rh",
			},
			wantErr: true,
			errType: "validation",
		},
		{
			name: "no results",
			input: RecommendInput{
				Input: "https://open.spotify.com/track/4iV5W9uYEdYUVa79Axb7Rh",
			},
			fetchRecommend: func(ctx context.Context, spotifyURL string, mode domain.RecommendMode, limit int) (*domain.RecommendResult, error) {
				return &domain.RecommendResult{
					SeedTrack: domain.Track{ID: "seed", Name: "Seed Track"},
					Items:     []domain.SimilarTrack{},
					Mode:      domain.RecommendModeBalanced,
				}, nil
			},
			wantErr: true,
			errType: "notfound",
		},
		{
			name: "repository error",
			input: RecommendInput{
				Input: "https://open.spotify.com/track/4iV5W9uYEdYUVa79Axb7Rh",
			},
			fetchRecommend: func(ctx context.Context, spotifyURL string, mode domain.RecommendMode, limit int) (*domain.RecommendResult, error) {
				return nil, errors.New("API error")
			},
			wantErr: true,
			errType: "other",
		},
		{
			name: "seed track fetch fails but recommend succeeds",
			input: RecommendInput{
				Input: "https://open.spotify.com/track/4iV5W9uYEdYUVa79Axb7Rh",
			},
			fetchRecommend: func(ctx context.Context, spotifyURL string, mode domain.RecommendMode, limit int) (*domain.RecommendResult, error) {
				return &domain.RecommendResult{
					SeedTrack: domain.Track{ID: "seed", Name: "Seed Track Basic"},
					Items: []domain.SimilarTrack{
						{ID: "1", Name: "Similar 1"},
					},
					Mode: domain.RecommendModeBalanced,
				}, nil
			},
			fetchTrack: func(ctx context.Context, spotifyURL string) (*domain.Track, error) {
				return nil, errors.New("fetch track failed")
			},
			wantErr:        false,
			wantItemsCount: 1,
		},
		{
			name: "with custom limit",
			input: RecommendInput{
				Input: "https://open.spotify.com/track/4iV5W9uYEdYUVa79Axb7Rh",
				Limit: 5,
			},
			fetchRecommend: func(ctx context.Context, spotifyURL string, mode domain.RecommendMode, limit int) (*domain.RecommendResult, error) {
				if limit != 5 {
					return nil, errors.New("wrong limit")
				}
				return &domain.RecommendResult{
					SeedTrack: domain.Track{ID: "seed", Name: "Seed Track"},
					Items: []domain.SimilarTrack{
						{ID: "1", Name: "Similar 1"},
					},
					Mode: domain.RecommendModeBalanced,
				}, nil
			},
			fetchTrack: func(ctx context.Context, spotifyURL string) (*domain.Track, error) {
				return &domain.Track{ID: "seed", Name: "Seed Track"}, nil
			},
			wantErr:        false,
			wantItemsCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockTrackRepository{
				fetchRecommendFunc: tt.fetchRecommend,
				fetchTrackFunc:     tt.fetchTrack,
			}
			uc := NewRecommendUseCase(repo)

			output, err := uc.GetRecommend(context.Background(), tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
					return
				}
				switch tt.errType {
				case "validation":
					if !IsValidationError(err) {
						t.Errorf("expected ValidationError but got %T: %v", err, err)
					}
				case "notfound":
					if !IsNotFoundError(err) {
						t.Errorf("expected NotFoundError but got %T: %v", err, err)
					}
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(output.Items) != tt.wantItemsCount {
				t.Errorf("Items count = %v, want %v", len(output.Items), tt.wantItemsCount)
			}
			if tt.wantMode != "" && output.Mode != tt.wantMode {
				t.Errorf("Mode = %v, want %v", output.Mode, tt.wantMode)
			}
		})
	}
}

func TestNewRecommendUseCase(t *testing.T) {
	repo := &mockTrackRepository{}
	uc := NewRecommendUseCase(repo)

	if uc == nil {
		t.Fatal("NewRecommendUseCase returned nil")
	}
	if uc.repo != repo {
		t.Error("repo not set correctly")
	}
}
