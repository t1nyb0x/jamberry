package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/t1nyb0x/jamberry/internal/domain"
)

// mockTrackRepository はTrackRepositoryのモック実装です
type mockTrackRepository struct {
	fetchTrackFunc     func(ctx context.Context, spotifyURL string) (*domain.Track, error)
	fetchSimilarFunc   func(ctx context.Context, spotifyURL string) ([]domain.SimilarTrack, error)
	fetchRecommendFunc func(ctx context.Context, spotifyURL string, mode domain.RecommendMode, limit int) (*domain.RecommendResult, error)
	searchTracksFunc   func(ctx context.Context, query string) ([]domain.Track, error)
}

func (m *mockTrackRepository) FetchTrack(ctx context.Context, spotifyURL string) (*domain.Track, error) {
	if m.fetchTrackFunc != nil {
		return m.fetchTrackFunc(ctx, spotifyURL)
	}
	return nil, errors.New("not implemented")
}

func (m *mockTrackRepository) FetchSimilar(ctx context.Context, spotifyURL string) ([]domain.SimilarTrack, error) {
	if m.fetchSimilarFunc != nil {
		return m.fetchSimilarFunc(ctx, spotifyURL)
	}
	return nil, errors.New("not implemented")
}

func (m *mockTrackRepository) FetchRecommend(ctx context.Context, spotifyURL string, mode domain.RecommendMode, limit int) (*domain.RecommendResult, error) {
	if m.fetchRecommendFunc != nil {
		return m.fetchRecommendFunc(ctx, spotifyURL, mode, limit)
	}
	return nil, errors.New("not implemented")
}

func (m *mockTrackRepository) SearchTracks(ctx context.Context, query string) ([]domain.Track, error) {
	if m.searchTracksFunc != nil {
		return m.searchTracksFunc(ctx, query)
	}
	return nil, errors.New("not implemented")
}

func TestTrackUseCase_GetTrack(t *testing.T) {
	tests := []struct {
		name      string
		input     TrackInput
		mockFunc  func(ctx context.Context, spotifyURL string) (*domain.Track, error)
		wantErr   bool
		errType   string
		wantTrack string
	}{
		{
			name: "valid URL",
			input: TrackInput{
				Input: "https://open.spotify.com/track/4iV5W9uYEdYUVa79Axb7Rh",
			},
			mockFunc: func(ctx context.Context, spotifyURL string) (*domain.Track, error) {
				return &domain.Track{
					ID:   "4iV5W9uYEdYUVa79Axb7Rh",
					Name: "Test Track",
				}, nil
			},
			wantErr:   false,
			wantTrack: "Test Track",
		},
		{
			name: "valid URI",
			input: TrackInput{
				Input: "spotify:track:4iV5W9uYEdYUVa79Axb7Rh",
			},
			mockFunc: func(ctx context.Context, spotifyURL string) (*domain.Track, error) {
				return &domain.Track{
					ID:   "4iV5W9uYEdYUVa79Axb7Rh",
					Name: "Test Track",
				}, nil
			},
			wantErr:   false,
			wantTrack: "Test Track",
		},
		{
			name: "valid ID",
			input: TrackInput{
				Input: "4iV5W9uYEdYUVa79Axb7Rh",
			},
			mockFunc: func(ctx context.Context, spotifyURL string) (*domain.Track, error) {
				return &domain.Track{
					ID:   "4iV5W9uYEdYUVa79Axb7Rh",
					Name: "Test Track",
				}, nil
			},
			wantErr:   false,
			wantTrack: "Test Track",
		},
		{
			name: "empty input",
			input: TrackInput{
				Input: "",
			},
			wantErr: true,
			errType: "validation",
		},
		{
			name: "invalid URL - not spotify",
			input: TrackInput{
				Input: "https://example.com/track/xxx",
			},
			wantErr: true,
			errType: "validation",
		},
		{
			name: "artist URL instead of track",
			input: TrackInput{
				Input: "https://open.spotify.com/artist/4iV5W9uYEdYUVa79Axb7Rh",
			},
			wantErr: true,
			errType: "validation",
		},
		{
			name: "repository error",
			input: TrackInput{
				Input: "https://open.spotify.com/track/4iV5W9uYEdYUVa79Axb7Rh",
			},
			mockFunc: func(ctx context.Context, spotifyURL string) (*domain.Track, error) {
				return nil, errors.New("API error")
			},
			wantErr: true,
			errType: "other",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockTrackRepository{
				fetchTrackFunc: tt.mockFunc,
			}
			uc := NewTrackUseCase(repo)

			output, err := uc.GetTrack(context.Background(), tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
					return
				}
				if tt.errType == "validation" && !IsValidationError(err) {
					t.Errorf("expected ValidationError but got %T", err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if output.Track.Name != tt.wantTrack {
				t.Errorf("Track.Name = %v, want %v", output.Track.Name, tt.wantTrack)
			}
		})
	}
}

func TestNewTrackUseCase(t *testing.T) {
	repo := &mockTrackRepository{}
	uc := NewTrackUseCase(repo)

	if uc == nil {
		t.Error("NewTrackUseCase returned nil")
	}
	if uc.repo != repo {
		t.Error("repo not set correctly")
	}
}
