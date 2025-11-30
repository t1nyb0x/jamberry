package usecase

import (
	"context"
	"log/slog"

	"github.com/t1nyb0x/jamberry/internal/domain"
	"github.com/t1nyb0x/jamberry/internal/spotify"
)

// ArtistUseCase はアーティスト関連のユースケースを提供します
type ArtistUseCase struct {
	repo domain.ArtistRepository
}

// NewArtistUseCase は新しいArtistUseCaseを作成します
func NewArtistUseCase(repo domain.ArtistRepository) *ArtistUseCase {
	return &ArtistUseCase{repo: repo}
}

// ArtistInput はアーティスト取得の入力パラメータです
type ArtistInput struct {
	Input string
}

// ArtistOutput はアーティスト取得の出力結果です
type ArtistOutput struct {
	Artist *domain.ArtistDetail
}

// GetArtist はアーティスト情報を取得します
func (u *ArtistUseCase) GetArtist(ctx context.Context, input ArtistInput) (*ArtistOutput, error) {
	// バリデーション
	result := spotify.ValidateInput(input.Input, spotify.EntityArtist)
	if !result.Valid {
		slog.Info("validation failed", "usecase", "artist", "input", input.Input, "error", result.Error)
		return nil, &ValidationError{Message: result.Error}
	}

	slog.Debug("validation passed", "usecase", "artist", "url", result.URL, "id", result.ID)

	// アーティスト情報を取得
	artist, err := u.repo.FetchArtist(ctx, result.URL)
	if err != nil {
		slog.Warn("artist fetch failed", "usecase", "artist", "url", result.URL, "error", err)
		return nil, err
	}

	slog.Info("artist fetched", "usecase", "artist", "artist_name", artist.Name, "artist_id", artist.ID)

	return &ArtistOutput{Artist: artist}, nil
}
