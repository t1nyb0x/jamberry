package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/t1nyb0x/jamberry/internal/domain"
)

// mockAlbumRepository はAlbumRepositoryのモック実装です
type mockAlbumRepository struct {
	fetchAlbumFunc func(ctx context.Context, spotifyURL string) (*domain.AlbumDetail, error)
}

func (m *mockAlbumRepository) FetchAlbum(ctx context.Context, spotifyURL string) (*domain.AlbumDetail, error) {
	if m.fetchAlbumFunc != nil {
		return m.fetchAlbumFunc(ctx, spotifyURL)
	}
	return nil, errors.New("not implemented")
}

func TestAlbumUseCase_GetAlbum(t *testing.T) {
	tests := []struct {
		name      string
		input     AlbumInput
		mockFunc  func(ctx context.Context, spotifyURL string) (*domain.AlbumDetail, error)
		wantErr   bool
		errType   string
		wantAlbum string
	}{
		{
			name: "valid URL",
			input: AlbumInput{
				Input: "https://open.spotify.com/album/0sNOF9WDwhWunNAHPD3Baj",
			},
			mockFunc: func(ctx context.Context, spotifyURL string) (*domain.AlbumDetail, error) {
				return &domain.AlbumDetail{
					ID:   "0sNOF9WDwhWunNAHPD3Baj",
					Name: "Test Album",
				}, nil
			},
			wantErr:   false,
			wantAlbum: "Test Album",
		},
		{
			name: "valid URI",
			input: AlbumInput{
				Input: "spotify:album:0sNOF9WDwhWunNAHPD3Baj",
			},
			mockFunc: func(ctx context.Context, spotifyURL string) (*domain.AlbumDetail, error) {
				return &domain.AlbumDetail{
					ID:   "0sNOF9WDwhWunNAHPD3Baj",
					Name: "Test Album",
				}, nil
			},
			wantErr:   false,
			wantAlbum: "Test Album",
		},
		{
			name: "valid ID",
			input: AlbumInput{
				Input: "0sNOF9WDwhWunNAHPD3Baj",
			},
			mockFunc: func(ctx context.Context, spotifyURL string) (*domain.AlbumDetail, error) {
				return &domain.AlbumDetail{
					ID:   "0sNOF9WDwhWunNAHPD3Baj",
					Name: "Test Album",
				}, nil
			},
			wantErr:   false,
			wantAlbum: "Test Album",
		},
		{
			name: "empty input",
			input: AlbumInput{
				Input: "",
			},
			wantErr: true,
			errType: "validation",
		},
		{
			name: "invalid URL - not spotify",
			input: AlbumInput{
				Input: "https://example.com/album/xxx",
			},
			wantErr: true,
			errType: "validation",
		},
		{
			name: "track URL instead of album",
			input: AlbumInput{
				Input: "https://open.spotify.com/track/4iV5W9uYEdYUVa79Axb7Rh",
			},
			wantErr: true,
			errType: "validation",
		},
		{
			name: "artist URL instead of album",
			input: AlbumInput{
				Input: "https://open.spotify.com/artist/0OdUWJ0sBjDrqHygGUXeCF",
			},
			wantErr: true,
			errType: "validation",
		},
		{
			name: "repository error",
			input: AlbumInput{
				Input: "https://open.spotify.com/album/0sNOF9WDwhWunNAHPD3Baj",
			},
			mockFunc: func(ctx context.Context, spotifyURL string) (*domain.AlbumDetail, error) {
				return nil, errors.New("API error")
			},
			wantErr: true,
			errType: "other",
		},
		{
			name: "URL with intl-ja",
			input: AlbumInput{
				Input: "https://open.spotify.com/intl-ja/album/0sNOF9WDwhWunNAHPD3Baj",
			},
			mockFunc: func(ctx context.Context, spotifyURL string) (*domain.AlbumDetail, error) {
				return &domain.AlbumDetail{
					ID:   "0sNOF9WDwhWunNAHPD3Baj",
					Name: "Test Album",
				}, nil
			},
			wantErr:   false,
			wantAlbum: "Test Album",
		},
		{
			name: "album with tracks",
			input: AlbumInput{
				Input: "https://open.spotify.com/album/0sNOF9WDwhWunNAHPD3Baj",
			},
			mockFunc: func(ctx context.Context, spotifyURL string) (*domain.AlbumDetail, error) {
				return &domain.AlbumDetail{
					ID:   "0sNOF9WDwhWunNAHPD3Baj",
					Name: "Test Album",
					Tracks: []domain.AlbumTrack{
						{ID: "track1", Name: "Track 1", TrackNumber: 1},
						{ID: "track2", Name: "Track 2", TrackNumber: 2},
					},
				}, nil
			},
			wantErr:   false,
			wantAlbum: "Test Album",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockAlbumRepository{
				fetchAlbumFunc: tt.mockFunc,
			}
			uc := NewAlbumUseCase(repo)

			output, err := uc.GetAlbum(context.Background(), tt.input)

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

			if output.Album.Name != tt.wantAlbum {
				t.Errorf("Album.Name = %v, want %v", output.Album.Name, tt.wantAlbum)
			}
		})
	}
}

func TestNewAlbumUseCase(t *testing.T) {
	repo := &mockAlbumRepository{}
	uc := NewAlbumUseCase(repo)

	if uc == nil {
		t.Error("NewAlbumUseCase returned nil")
	}
	if uc.repo != repo {
		t.Error("repo not set correctly")
	}
}
