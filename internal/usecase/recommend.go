package usecase

import (
	"context"
	"log/slog"

	"github.com/t1nyb0x/jamberry/internal/domain"
	"github.com/t1nyb0x/jamberry/internal/spotify"
)

// RecommendUseCase はレコメンド関連のユースケースを提供します
type RecommendUseCase struct {
	repo domain.TrackRepository
}

// NewRecommendUseCase は新しいRecommendUseCaseを作成します
func NewRecommendUseCase(repo domain.TrackRepository) *RecommendUseCase {
	return &RecommendUseCase{repo: repo}
}

// RecommendInput はレコメンド取得の入力パラメータです
type RecommendInput struct {
	Input string
}

// RecommendOutput はレコメンド取得の出力結果です
type RecommendOutput struct {
	SourceTrack   *domain.Track
	SimilarTracks []domain.SimilarTrack
}

// GetRecommend はレコメンド情報を取得します
func (u *RecommendUseCase) GetRecommend(ctx context.Context, input RecommendInput) (*RecommendOutput, error) {
	// バリデーション（トラックURLのみ受け付ける）
	result := spotify.ValidateInput(input.Input, spotify.EntityTrack)
	if !result.Valid {
		slog.Info("validation failed", "usecase", "recommend", "input", input.Input, "error", result.Error)
		return nil, &ValidationError{Message: result.Error}
	}

	slog.Debug("validation passed", "usecase", "recommend", "url", result.URL, "id", result.ID)

	// 元のトラック情報を取得
	track, err := u.repo.FetchTrack(ctx, result.URL)
	if err != nil {
		slog.Warn("track fetch failed for recommend", "usecase", "recommend", "url", result.URL, "error", err)
		return nil, err
	}

	// 類似トラックを取得
	similar, err := u.repo.FetchSimilar(ctx, result.URL)
	if err != nil {
		slog.Warn("similar fetch failed", "usecase", "recommend", "url", result.URL, "error", err)
		return nil, err
	}

	if len(similar) == 0 {
		slog.Info("no similar tracks found", "usecase", "recommend", "track_name", track.Name)
		return nil, &NotFoundError{Message: "🔍 該当する結果は見つかりませんでした。"}
	}

	slog.Info("recommend fetched", "usecase", "recommend", "track_name", track.Name, "result_count", len(similar))

	return &RecommendOutput{
		SourceTrack:   track,
		SimilarTracks: similar,
	}, nil
}
