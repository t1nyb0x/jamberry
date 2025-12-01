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
	Mode  domain.RecommendMode
	Limit int
}

// RecommendOutput はレコメンド取得の出力結果です
type RecommendOutput struct {
	SeedTrack    *domain.Track
	SeedFeatures *domain.TrackFeatures // v2: Deezer + MusicBrainz features
	Items        []domain.SimilarTrack
	Mode         domain.RecommendMode
}

// GetRecommend はレコメンド情報を取得します
func (u *RecommendUseCase) GetRecommend(ctx context.Context, input RecommendInput) (*RecommendOutput, error) {
	// バリデーション（トラックURLのみ受け付ける）
	result := spotify.ValidateInput(input.Input, spotify.EntityTrack)
	if !result.Valid {
		slog.Info("validation failed", "usecase", "recommend", "input", input.Input, "error", result.Error)
		return nil, &ValidationError{Message: result.Error}
	}

	slog.Debug("validation passed", "usecase", "recommend", "url", result.URL, "id", result.ID, "mode", input.Mode)

	// デフォルト値の設定
	mode := input.Mode
	if mode == "" {
		mode = domain.RecommendModeBalanced
	}
	limit := input.Limit
	if limit <= 0 {
		limit = 20
	}

	// 新しいレコメンドAPIを使用
	recommendResult, err := u.repo.FetchRecommend(ctx, result.URL, mode, limit)
	if err != nil {
		slog.Warn("recommend fetch failed", "usecase", "recommend", "url", result.URL, "mode", mode, "error", err)
		return nil, err
	}

	if len(recommendResult.Items) == 0 {
		slog.Info("no recommend tracks found", "usecase", "recommend", "track_name", recommendResult.SeedTrack.Name)
		return nil, &NotFoundError{Message: "🔍 該当する結果は見つかりませんでした。"}
	}

	slog.Info("recommend fetched", "usecase", "recommend",
		"track_name", recommendResult.SeedTrack.Name,
		"mode", recommendResult.Mode,
		"result_count", len(recommendResult.Items))

	// SeedTrackの完全な情報を取得（URLやDurationMsなどが必要な場合）
	seedTrack, err := u.repo.FetchTrack(ctx, result.URL)
	if err != nil {
		// シードトラックの詳細取得に失敗しても、レコメンド結果は返す
		slog.Warn("seed track detail fetch failed, using basic info", "url", result.URL, "error", err)
		seedTrack = &recommendResult.SeedTrack
	}

	return &RecommendOutput{
		SeedTrack:    seedTrack,
		SeedFeatures: recommendResult.SeedFeatures,
		Items:        recommendResult.Items,
		Mode:         recommendResult.Mode,
	}, nil
}
