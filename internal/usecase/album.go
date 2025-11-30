package usecase

import (
	"context"
	"log/slog"

	"github.com/t1nyb0x/jamberry/internal/domain"
	"github.com/t1nyb0x/jamberry/internal/spotify"
)

// AlbumUseCase はアルバム関連のユースケースを提供します
type AlbumUseCase struct {
	repo domain.AlbumRepository
}

// NewAlbumUseCase は新しいAlbumUseCaseを作成します
func NewAlbumUseCase(repo domain.AlbumRepository) *AlbumUseCase {
	return &AlbumUseCase{repo: repo}
}

// AlbumInput はアルバム取得の入力パラメータです
type AlbumInput struct {
	Input string
}

// AlbumOutput はアルバム取得の出力結果です
type AlbumOutput struct {
	Album *domain.AlbumDetail
}

// GetAlbum はアルバム情報を取得します
func (u *AlbumUseCase) GetAlbum(ctx context.Context, input AlbumInput) (*AlbumOutput, error) {
	// バリデーション
	result := spotify.ValidateInput(input.Input, spotify.EntityAlbum)
	if !result.Valid {
		slog.Info("validation failed", "usecase", "album", "input", input.Input, "error", result.Error)
		return nil, &ValidationError{Message: result.Error}
	}

	slog.Debug("validation passed", "usecase", "album", "url", result.URL, "id", result.ID)

	// アルバム情報を取得
	album, err := u.repo.FetchAlbum(ctx, result.URL)
	if err != nil {
		slog.Warn("album fetch failed", "usecase", "album", "url", result.URL, "error", err)
		return nil, err
	}

	slog.Info("album fetched", "usecase", "album", "album_name", album.Name, "album_id", album.ID)

	return &AlbumOutput{Album: album}, nil
}
