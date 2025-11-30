package usecase

import (
	"context"
	"log/slog"

	"github.com/t1nyb0x/jamberry/internal/domain"
	"github.com/t1nyb0x/jamberry/internal/spotify"
)

// TrackUseCase はトラック関連のユースケースを提供します
type TrackUseCase struct {
	repo domain.TrackRepository
}

// NewTrackUseCase は新しいTrackUseCaseを作成します
func NewTrackUseCase(repo domain.TrackRepository) *TrackUseCase {
	return &TrackUseCase{repo: repo}
}

// TrackInput はトラック取得の入力パラメータです
type TrackInput struct {
	Input string
}

// TrackOutput はトラック取得の出力結果です
type TrackOutput struct {
	Track *domain.Track
}

// GetTrack はトラック情報を取得します
func (u *TrackUseCase) GetTrack(ctx context.Context, input TrackInput) (*TrackOutput, error) {
	// バリデーション
	result := spotify.ValidateInput(input.Input, spotify.EntityTrack)
	if !result.Valid {
		slog.Info("validation failed", "usecase", "track", "input", input.Input, "error", result.Error)
		return nil, &ValidationError{Message: result.Error}
	}

	slog.Debug("validation passed", "usecase", "track", "url", result.URL, "id", result.ID)

	// トラック情報を取得
	track, err := u.repo.FetchTrack(ctx, result.URL)
	if err != nil {
		slog.Warn("track fetch failed", "usecase", "track", "url", result.URL, "error", err)
		return nil, err
	}

	slog.Info("track fetched", "usecase", "track", "track_name", track.Name, "track_id", track.ID)

	return &TrackOutput{Track: track}, nil
}
