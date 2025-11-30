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
)

// Client はtracktaste APIクライアントです
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient は新しいtracktaste APIクライアントを作成します
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// APIError はtracktaste APIからのエラーレスポンスを表します
type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("tracktaste API error: %s (code: %s, status: %d)", e.Message, e.Code, e.Status)
}

// Image は画像情報を表します
type Image struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

// Artist はアーティスト情報を表します
type Artist struct {
	URL  string `json:"url"`
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Album はアルバム情報を表します
type Album struct {
	URL         string   `json:"url"`
	ID          string   `json:"id"`
	Images      []Image  `json:"images"`
	Name        string   `json:"name"`
	ReleaseDate string   `json:"release_date"`
	Artists     []Artist `json:"artists"`
}

// Track はトラック情報を表します
type Track struct {
	Album       Album    `json:"album"`
	Artists     []Artist `json:"artists"`
	DiscNumber  int      `json:"disc_number"`
	Popularity  *int     `json:"popularity"`
	ISRC        *string  `json:"isrc"`
	URL         string   `json:"url"`
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	TrackNumber int      `json:"track_number"`
	DurationMs  int      `json:"duration_ms"`
	Explicit    bool     `json:"explicit"`
}

// SimilarTrack は類似トラック情報を表します
type SimilarTrack struct {
	Album       Album   `json:"album"`
	ISRC        *string `json:"isrc"`
	UPC         *string `json:"upc"`
	URL         string  `json:"url"`
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Popularity  *int    `json:"popularity"`
	TrackNumber int     `json:"track_number"`
	DurationMs  int     `json:"duration_ms"`
	Explicit    bool    `json:"explicit"`
}

// SearchTrack は検索結果のトラック情報を表します
type SearchTrack struct {
	Album       Album    `json:"album"`
	Artists     []Artist `json:"artists"`
	DiscNumber  int      `json:"disc_number"`
	Popularity  *int     `json:"popularity"`
	ISRC        *string  `json:"isrc"`
	URL         string   `json:"url"`
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	TrackNumber int      `json:"track_number"`
	DurationMs  int      `json:"duration_ms"`
	Explicit    bool     `json:"explicit"`
}

// ArtistFull はアーティストの詳細情報を表します
type ArtistFull struct {
	URL        string   `json:"url"`
	Followers  string   `json:"followers"`
	Genres     []string `json:"genres"`
	ID         string   `json:"id"`
	Images     []Image  `json:"images"`
	Name       string   `json:"name"`
	Popularity *int     `json:"popularity"`
}

// AlbumTrack はアルバム内のトラック情報を表します
type AlbumTrack struct {
	Artists     []Artist `json:"artists"`
	URL         string   `json:"url"`
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	TrackNumber int      `json:"track_number"`
}

// AlbumTracks はアルバム内のトラックリストを表します
type AlbumTracks struct {
	Items []AlbumTrack `json:"items"`
}

// AlbumFull はアルバムの詳細情報を表します
type AlbumFull struct {
	URL         string      `json:"url"`
	ID          string      `json:"id"`
	Images      []Image     `json:"images"`
	Name        string      `json:"name"`
	ReleaseDate string      `json:"release_date"`
	Artists     []Artist    `json:"artists"`
	Tracks      AlbumTracks `json:"tracks"`
	Popularity  *int        `json:"popularity"`
	UPC         string      `json:"upc"`
	Genres      []string    `json:"genres"`
}

// Response はtracktaste APIの共通レスポンス形式を表します
type Response[T any] struct {
	Status int `json:"status"`
	Result T   `json:"result"`
}

// SimilarResponse は類似トラックのレスポンス形式を表します
type SimilarResponse struct {
	Items []SimilarTrack `json:"items"`
}

// SearchResponse は検索結果のレスポンス形式を表します
type SearchResponse struct {
	Items []SearchTrack `json:"items"`
}

// FetchTrack はトラック情報を取得します
func (c *Client) FetchTrack(ctx context.Context, spotifyURL string) (*Track, error) {
	endpoint := fmt.Sprintf("%s/v1/track/fetch?url=%s", c.baseURL, url.QueryEscape(spotifyURL))
	return doRequest[Track](ctx, c, endpoint)
}

// FetchSimilar は類似トラックを取得します
func (c *Client) FetchSimilar(ctx context.Context, spotifyURL string) ([]SimilarTrack, error) {
	endpoint := fmt.Sprintf("%s/v1/track/similar?url=%s", c.baseURL, url.QueryEscape(spotifyURL))
	resp, err := doRequest[SimilarResponse](ctx, c, endpoint)
	if err != nil {
		return nil, err
	}
	return resp.Items, nil
}

// SearchTracks はトラックを検索します
func (c *Client) SearchTracks(ctx context.Context, query string) ([]SearchTrack, error) {
	endpoint := fmt.Sprintf("%s/v1/track/search?q=%s", c.baseURL, url.QueryEscape(query))
	resp, err := doRequest[SearchResponse](ctx, c, endpoint)
	if err != nil {
		return nil, err
	}
	return resp.Items, nil
}

// FetchArtist はアーティスト情報を取得します
func (c *Client) FetchArtist(ctx context.Context, spotifyURL string) (*ArtistFull, error) {
	endpoint := fmt.Sprintf("%s/v1/artist/fetch?url=%s", c.baseURL, url.QueryEscape(spotifyURL))
	return doRequest[ArtistFull](ctx, c, endpoint)
}

// FetchAlbum はアルバム情報を取得します
func (c *Client) FetchAlbum(ctx context.Context, spotifyURL string) (*AlbumFull, error) {
	endpoint := fmt.Sprintf("%s/v1/album/fetch?url=%s", c.baseURL, url.QueryEscape(spotifyURL))
	return doRequest[AlbumFull](ctx, c, endpoint)
}

// doRequest はAPIリクエストを実行します
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
			// JSONパースに失敗した場合はステータスコードに応じたエラーを返す
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
