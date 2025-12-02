package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/t1nyb0x/jamberry/internal/domain"
)

// mockArtistRepository はArtistRepositoryのモック実装です
type mockArtistRepository struct {
	fetchArtistFunc func(ctx context.Context, spotifyURL string) (*domain.ArtistDetail, error)
}

func (m *mockArtistRepository) FetchArtist(ctx context.Context, spotifyURL string) (*domain.ArtistDetail, error) {
	if m.fetchArtistFunc != nil {
		return m.fetchArtistFunc(ctx, spotifyURL)
	}
	return nil, errors.New("not implemented")
}

func TestArtistUseCase_GetArtist(t *testing.T) {
	tests := []struct {
		name       string
		input      ArtistInput
		mockFunc   func(ctx context.Context, spotifyURL string) (*domain.ArtistDetail, error)
		wantErr    bool
		errType    string
		wantArtist string
	}{
		{
			name: "valid URL",
			input: ArtistInput{
				Input: "https://open.spotify.com/artist/0OdUWJ0sBjDrqHygGUXeCF",
			},
			mockFunc: func(ctx context.Context, spotifyURL string) (*domain.ArtistDetail, error) {
				return &domain.ArtistDetail{
					ID:   "0OdUWJ0sBjDrqHygGUXeCF",
					Name: "Test Artist",
				}, nil
			},
			wantErr:    false,
			wantArtist: "Test Artist",
		},
		{
			name: "valid URI",
			input: ArtistInput{
				Input: "spotify:artist:0OdUWJ0sBjDrqHygGUXeCF",
			},
			mockFunc: func(ctx context.Context, spotifyURL string) (*domain.ArtistDetail, error) {
				return &domain.ArtistDetail{
					ID:   "0OdUWJ0sBjDrqHygGUXeCF",
					Name: "Test Artist",
				}, nil
			},
			wantErr:    false,
			wantArtist: "Test Artist",
		},
		{
			name: "valid ID",
			input: ArtistInput{
				Input: "0OdUWJ0sBjDrqHygGUXeCF",
			},
			mockFunc: func(ctx context.Context, spotifyURL string) (*domain.ArtistDetail, error) {
				return &domain.ArtistDetail{
					ID:   "0OdUWJ0sBjDrqHygGUXeCF",
					Name: "Test Artist",
				}, nil
			},
			wantErr:    false,
			wantArtist: "Test Artist",
		},
		{
			name: "empty input",
			input: ArtistInput{
				Input: "",
			},
			wantErr: true,
			errType: "validation",
		},
		{
			name: "invalid URL - not spotify",
			input: ArtistInput{
				Input: "https://example.com/artist/xxx",
			},
			wantErr: true,
			errType: "validation",
		},
		{
			name: "track URL instead of artist",
			input: ArtistInput{
				Input: "https://open.spotify.com/track/4iV5W9uYEdYUVa79Axb7Rh",
			},
			wantErr: true,
			errType: "validation",
		},
		{
			name: "album URL instead of artist",
			input: ArtistInput{
				Input: "https://open.spotify.com/album/4iV5W9uYEdYUVa79Axb7Rh",
			},
			wantErr: true,
			errType: "validation",
		},
		{
			name: "repository error",
			input: ArtistInput{
				Input: "https://open.spotify.com/artist/0OdUWJ0sBjDrqHygGUXeCF",
			},
			mockFunc: func(ctx context.Context, spotifyURL string) (*domain.ArtistDetail, error) {
				return nil, errors.New("API error")
			},
			wantErr: true,
			errType: "other",
		},
		{
			name: "URL with intl-ja",
			input: ArtistInput{
				Input: "https://open.spotify.com/intl-ja/artist/0OdUWJ0sBjDrqHygGUXeCF",
			},
			mockFunc: func(ctx context.Context, spotifyURL string) (*domain.ArtistDetail, error) {
				return &domain.ArtistDetail{
					ID:   "0OdUWJ0sBjDrqHygGUXeCF",
					Name: "Test Artist",
				}, nil
			},
			wantErr:    false,
			wantArtist: "Test Artist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockArtistRepository{
				fetchArtistFunc: tt.mockFunc,
			}
			uc := NewArtistUseCase(repo)

			output, err := uc.GetArtist(context.Background(), tt.input)

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

			if output.Artist.Name != tt.wantArtist {
				t.Errorf("Artist.Name = %v, want %v", output.Artist.Name, tt.wantArtist)
			}
		})
	}
}

func TestNewArtistUseCase(t *testing.T) {
	repo := &mockArtistRepository{}
	uc := NewArtistUseCase(repo)

	if uc == nil {
		t.Fatal("NewArtistUseCase returned nil")
	}
	if uc.repo != repo {
		t.Error("repo not set correctly")
	}
}
