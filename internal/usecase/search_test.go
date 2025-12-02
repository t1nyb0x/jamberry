package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/t1nyb0x/jamberry/internal/domain"
)

func TestSearchUseCase_SearchTracks(t *testing.T) {
	tests := []struct {
		name       string
		input      SearchInput
		mockFunc   func(ctx context.Context, query string) ([]domain.Track, error)
		wantErr    bool
		errType    string
		wantCount  int
		wantQuery  string
	}{
		{
			name: "valid query with results",
			input: SearchInput{
				Query: "test query",
			},
			mockFunc: func(ctx context.Context, query string) ([]domain.Track, error) {
				return []domain.Track{
					{ID: "1", Name: "Track 1"},
					{ID: "2", Name: "Track 2"},
				}, nil
			},
			wantErr:   false,
			wantCount: 2,
			wantQuery: "test query",
		},
		{
			name: "query with leading/trailing spaces",
			input: SearchInput{
				Query: "  trimmed query  ",
			},
			mockFunc: func(ctx context.Context, query string) ([]domain.Track, error) {
				// スペースがトリムされていることを確認
				if query != "trimmed query" {
					return nil, errors.New("query was not trimmed")
				}
				return []domain.Track{
					{ID: "1", Name: "Track 1"},
				}, nil
			},
			wantErr:   false,
			wantCount: 1,
			wantQuery: "trimmed query",
		},
		{
			name: "empty query",
			input: SearchInput{
				Query: "",
			},
			wantErr: true,
			errType: "validation",
		},
		{
			name: "whitespace only query",
			input: SearchInput{
				Query: "   ",
			},
			wantErr: true,
			errType: "validation",
		},
		{
			name: "no results",
			input: SearchInput{
				Query: "nonexistent track",
			},
			mockFunc: func(ctx context.Context, query string) ([]domain.Track, error) {
				return []domain.Track{}, nil
			},
			wantErr: true,
			errType: "notfound",
		},
		{
			name: "repository error",
			input: SearchInput{
				Query: "test query",
			},
			mockFunc: func(ctx context.Context, query string) ([]domain.Track, error) {
				return nil, errors.New("API error")
			},
			wantErr: true,
			errType: "other",
		},
		{
			name: "japanese query",
			input: SearchInput{
				Query: "米津玄師 Lemon",
			},
			mockFunc: func(ctx context.Context, query string) ([]domain.Track, error) {
				return []domain.Track{
					{ID: "1", Name: "Lemon"},
				}, nil
			},
			wantErr:   false,
			wantCount: 1,
			wantQuery: "米津玄師 Lemon",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockTrackRepository{
				searchTracksFunc: tt.mockFunc,
			}
			uc := NewSearchUseCase(repo)

			output, err := uc.SearchTracks(context.Background(), tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
					return
				}
				switch tt.errType {
				case "validation":
					if !IsValidationError(err) {
						t.Errorf("expected ValidationError but got %T", err)
					}
				case "notfound":
					if !IsNotFoundError(err) {
						t.Errorf("expected NotFoundError but got %T", err)
					}
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(output.Tracks) != tt.wantCount {
				t.Errorf("Tracks count = %v, want %v", len(output.Tracks), tt.wantCount)
			}
			if output.Query != tt.wantQuery {
				t.Errorf("Query = %v, want %v", output.Query, tt.wantQuery)
			}
		})
	}
}

func TestNewSearchUseCase(t *testing.T) {
	repo := &mockTrackRepository{}
	uc := NewSearchUseCase(repo)

	if uc == nil {
		t.Error("NewSearchUseCase returned nil")
	}
	if uc.repo != repo {
		t.Error("repo not set correctly")
	}
}
