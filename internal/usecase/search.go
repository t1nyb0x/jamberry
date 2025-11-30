package usecase

import (
	"context"
	"log/slog"
	"strings"

	"github.com/t1nyb0x/jamberry/internal/domain"
)

// SearchUseCase は検索関連のユースケースを提供します
type SearchUseCase struct {
	repo domain.TrackRepository
}

// NewSearchUseCase は新しいSearchUseCaseを作成します
func NewSearchUseCase(repo domain.TrackRepository) *SearchUseCase {
	return &SearchUseCase{repo: repo}
}

// SearchInput は検索の入力パラメータです
type SearchInput struct {
	Query string
}

// SearchOutput は検索の出力結果です
type SearchOutput struct {
	Query  string
	Tracks []domain.Track
}

// SearchTracks はトラックを検索します
func (u *SearchUseCase) SearchTracks(ctx context.Context, input SearchInput) (*SearchOutput, error) {
	query := strings.TrimSpace(input.Query)
	if query == "" {
		slog.Info("validation failed: empty query", "usecase", "search")
		return nil, &ValidationError{Message: "❌ 検索キーワードを入力してください。"}
	}

	slog.Debug("search query received", "usecase", "search", "query", query)

	// 検索実行
	tracks, err := u.repo.SearchTracks(ctx, query)
	if err != nil {
		slog.Warn("search failed", "usecase", "search", "query", query, "error", err)
		return nil, err
	}

	if len(tracks) == 0 {
		slog.Info("no results found", "usecase", "search", "query", query)
		return nil, &NotFoundError{Message: "🔍 該当する結果は見つかりませんでした。"}
	}

	slog.Info("search completed", "usecase", "search", "query", query, "result_count", len(tracks))

	return &SearchOutput{
		Query:  query,
		Tracks: tracks,
	}, nil
}
