package tracktaste

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/t1nyb0x/jamberry/internal/domain"
)

// Client はtracktaste APIクライアントです
// domain.MusicRepository インターフェースを実装します
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient は新しいtracktaste APIクライアントを作成します
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second, // v2 recommend APIは複数の外部APIを呼び出すため長めに設定
		},
	}
}

// インターフェース実装の確認
var _ domain.MusicRepository = (*Client)(nil)

// APIError はtracktaste APIからのエラーレスポンスを表します
type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("tracktaste API error: %s (code: %s, status: %d)", e.Message, e.Code, e.Status)
}

// Response はtracktaste APIの共通レスポンス形式を表します
type Response[T any] struct {
	Status int `json:"status"`
	Result T   `json:"result"`
}

// FetchTrack はトラック情報を取得します
func (c *Client) FetchTrack(ctx context.Context, spotifyURL string) (*domain.Track, error) {
	endpoint := fmt.Sprintf("%s/v1/track/fetch?url=%s", c.baseURL, url.QueryEscape(spotifyURL))

	resp, err := doRequest[trackResponse](ctx, c, endpoint)
	if err != nil {
		return nil, err
	}

	return resp.toDomain(), nil
}

// FetchSimilar は類似トラックを取得します（従来API互換）
func (c *Client) FetchSimilar(ctx context.Context, spotifyURL string) ([]domain.SimilarTrack, error) {
	endpoint := fmt.Sprintf("%s/v1/track/similar?url=%s", c.baseURL, url.QueryEscape(spotifyURL))

	resp, err := doRequest[similarResponse](ctx, c, endpoint)
	if err != nil {
		return nil, err
	}

	tracks := make([]domain.SimilarTrack, len(resp.Items))
	for i, item := range resp.Items {
		tracks[i] = item.toDomain()
	}

	return tracks, nil
}

// FetchRecommend はレコメンドトラックを取得します（v2 API: Deezer + MusicBrainz）
func (c *Client) FetchRecommend(ctx context.Context, spotifyURL string, mode domain.RecommendMode, limit int) (*domain.RecommendResult, error) {
	// デフォルト値の設定
	if mode == "" {
		mode = domain.RecommendModeBalanced
	}
	if limit <= 0 {
		limit = 20 // TrackTaste v2 APIのデフォルトは20件
	}
	if limit > 50 {
		limit = 50 // TrackTaste v2 APIの上限は50件
	}

	// v2 API エンドポイントを使用（レスポンス形式はv1と同じ: status + result）
	endpoint := fmt.Sprintf("%s/v2/track/recommend?url=%s&mode=%s&limit=%d",
		c.baseURL, url.QueryEscape(spotifyURL), string(mode), limit)

	resp, err := doRequest[recommendResponseV2](ctx, c, endpoint)
	if err != nil {
		return nil, err
	}

	return resp.toDomain(), nil
}

// SearchTracks はトラックを検索します
func (c *Client) SearchTracks(ctx context.Context, query string) ([]domain.Track, error) {
	endpoint := fmt.Sprintf("%s/v1/track/search?q=%s", c.baseURL, url.QueryEscape(query))

	resp, err := doRequest[searchResponse](ctx, c, endpoint)
	if err != nil {
		return nil, err
	}

	tracks := make([]domain.Track, len(resp.Items))
	for i, item := range resp.Items {
		tracks[i] = item.toDomain()
	}

	return tracks, nil
}

// FetchArtist はアーティスト情報を取得します
func (c *Client) FetchArtist(ctx context.Context, spotifyURL string) (*domain.ArtistDetail, error) {
	endpoint := fmt.Sprintf("%s/v1/artist/fetch?url=%s", c.baseURL, url.QueryEscape(spotifyURL))

	resp, err := doRequest[artistResponse](ctx, c, endpoint)
	if err != nil {
		return nil, err
	}

	return resp.toDomain(), nil
}

// FetchAlbum はアルバム情報を取得します
func (c *Client) FetchAlbum(ctx context.Context, spotifyURL string) (*domain.AlbumDetail, error) {
	endpoint := fmt.Sprintf("%s/v1/album/fetch?url=%s", c.baseURL, url.QueryEscape(spotifyURL))

	resp, err := doRequest[albumResponse](ctx, c, endpoint)
	if err != nil {
		return nil, err
	}

	return resp.toDomain(), nil
}

// doRequest はAPIリクエストを実行します（v1 API用）
func doRequest[T any](ctx context.Context, c *Client, endpoint string) (*T, error) {
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		slog.Warn("tracktaste API request failed",
			"endpoint", endpoint,
			"error", err,
			"latency_ms", time.Since(start).Milliseconds(),
		)
		return nil, fmt.Errorf("❌ 接続エラーが発生しました。")
	}
	defer resp.Body.Close()

	slog.Debug("tracktaste API request",
		"endpoint", endpoint,
		"status", resp.StatusCode,
		"latency_ms", time.Since(start).Milliseconds(),
	)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// エラーレスポンスの場合
	if resp.StatusCode != http.StatusOK {
		var apiErr APIError
		if err := json.Unmarshal(body, &apiErr); err != nil {
			return nil, handleHTTPError(resp.StatusCode)
		}
		return nil, handleAPIError(&apiErr)
	}

	var result Response[T]
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result.Result, nil
}

// handleHTTPError はHTTPステータスコードに応じたエラーを返します
func handleHTTPError(statusCode int) error {
	switch statusCode {
	case http.StatusServiceUnavailable:
		slog.Error("tracktaste API service unavailable", "status", statusCode)
		return fmt.Errorf("❌ サーバーエラーが発生しました。しばらくしてから再試行してください。")
	case http.StatusGatewayTimeout:
		slog.Warn("tracktaste API timeout", "status", statusCode)
		return fmt.Errorf("❌ リクエストがタイムアウトしました。しばらくしてから再試行してください。")
	case http.StatusTooManyRequests:
		slog.Warn("tracktaste API rate limited", "status", statusCode)
		return fmt.Errorf("⏳ リクエスト制限中です。しばらくしてから再試行してください。")
	default:
		slog.Error("tracktaste API error", "status", statusCode)
		return fmt.Errorf("❌ サーバーエラーが発生しました。しばらくしてから再試行してください。")
	}
}

// handleAPIError はAPIエラーコードに応じたエラーを返します
func handleAPIError(apiErr *APIError) error {
	switch apiErr.Code {
	case "EMPTY_PARAM", "EMPTY_URL":
		return fmt.Errorf("❌ URL を入力してください。")
	case "EMPTY_QUERY":
		return fmt.Errorf("❌ 検索キーワードを入力してください。")
	case "INVALID_PARAM":
		return fmt.Errorf("❌ 入力形式が不正です。")
	case "NOT_SPOTIFY_URL", "INVALID_URL":
		return fmt.Errorf("❌ Spotify の URL を入力してください。")
	case "INVALID_RESOURCE_TYPE", "DIFFERENT_SPOTIFY_URL":
		return fmt.Errorf("❌ 正しい種類の URL を入力してください。")
	case "SOMETHING_SPOTIFY_ERROR", "SOMETHING_KKBOX_ERROR":
		slog.Error("tracktaste API backend error", "code", apiErr.Code, "message", apiErr.Message)
		return fmt.Errorf("❌ サーバーエラーが発生しました。しばらくしてから再試行してください。")
	case "REQUEST_TIMEOUT":
		slog.Warn("tracktaste API timeout", "code", apiErr.Code)
		return fmt.Errorf("❌ リクエストがタイムアウトしました。しばらくしてから再試行してください。")
	default:
		slog.Error("tracktaste API unknown error", "code", apiErr.Code, "message", apiErr.Message)
		return fmt.Errorf("❌ サーバーエラーが発生しました。しばらくしてから再試行してください。")
	}
}
