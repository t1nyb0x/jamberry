package spotify

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// EntityType はSpotifyのエンティティ種別を表します
type EntityType string

const (
	EntityTrack   EntityType = "track"
	EntityArtist  EntityType = "artist"
	EntityAlbum   EntityType = "album"
	EntityUnknown EntityType = "unknown"
)

// SpotifyIDRegex はSpotify IDの形式を検証する正規表現
var SpotifyIDRegex = regexp.MustCompile(`^[a-zA-Z0-9]{22}$`)

// ValidationResult はバリデーション結果を表します
type ValidationResult struct {
	Valid      bool
	URL        string
	ID         string
	EntityType EntityType
	Error      string
}

// ValidateInput は入力をバリデーションし、Spotify URLに正規化します
func ValidateInput(input string, expectedType EntityType) ValidationResult {
	input = strings.TrimSpace(input)

	if input == "" {
		return ValidationResult{
			Valid: false,
			Error: "❌ Spotify の URL / ID として認識できませんでした。",
		}
	}

	// URL形式の場合
	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		return validateURL(input, expectedType)
	}

	// URI形式の場合 (spotify:track:xxx)
	if strings.HasPrefix(input, "spotify:") {
		return validateURI(input, expectedType)
	}

	// 生IDの場合
	return validateID(input, expectedType)
}

// validateURL はURL形式の入力をバリデーションします
func validateURL(input string, expectedType EntityType) ValidationResult {
	parsed, err := url.Parse(input)
	if err != nil {
		return ValidationResult{
			Valid: false,
			Error: "❌ Spotify の URL / ID として認識できませんでした。",
		}
	}

	// ドメインチェック
	if parsed.Host != "open.spotify.com" {
		return ValidationResult{
			Valid: false,
			Error: "❌ Spotify の URL / ID として認識できませんでした。",
		}
	}

	// パスからエンティティ種別とIDを抽出
	parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	if len(parts) < 2 {
		return ValidationResult{
			Valid: false,
			Error: "❌ Spotify の URL / ID として認識できませんでした。",
		}
	}

	entityType := EntityType(parts[0])
	id := parts[1]

	// クエリパラメータを除去
	if idx := strings.Index(id, "?"); idx != -1 {
		id = id[:idx]
	}

	// エンティティ種別チェック
	if entityType != expectedType {
		return ValidationResult{
			Valid: false,
			Error: getEntityMismatchError(expectedType),
		}
	}

	// IDの形式チェック
	if !SpotifyIDRegex.MatchString(id) {
		return ValidationResult{
			Valid: false,
			Error: "❌ Spotify の URL / ID として認識できませんでした。",
		}
	}

	return ValidationResult{
		Valid:      true,
		URL:        fmt.Sprintf("https://open.spotify.com/%s/%s", entityType, id),
		ID:         id,
		EntityType: entityType,
	}
}

// validateURI はURI形式の入力をバリデーションします
func validateURI(input string, expectedType EntityType) ValidationResult {
	parts := strings.Split(input, ":")
	if len(parts) != 3 {
		return ValidationResult{
			Valid: false,
			Error: "❌ Spotify の URL / ID として認識できませんでした。",
		}
	}

	entityType := EntityType(parts[1])
	id := parts[2]

	// エンティティ種別チェック
	if entityType != expectedType {
		return ValidationResult{
			Valid: false,
			Error: getEntityMismatchError(expectedType),
		}
	}

	// IDの形式チェック
	if !SpotifyIDRegex.MatchString(id) {
		return ValidationResult{
			Valid: false,
			Error: "❌ Spotify の URL / ID として認識できませんでした。",
		}
	}

	return ValidationResult{
		Valid:      true,
		URL:        fmt.Sprintf("https://open.spotify.com/%s/%s", entityType, id),
		ID:         id,
		EntityType: entityType,
	}
}

// validateID は生ID形式の入力をバリデーションします
func validateID(input string, expectedType EntityType) ValidationResult {
	// IDの形式チェック
	if !SpotifyIDRegex.MatchString(input) {
		return ValidationResult{
			Valid: false,
			Error: "❌ Spotify の URL / ID として認識できませんでした。",
		}
	}

	// 生IDの場合はエンティティ種別の検証は行わない
	return ValidationResult{
		Valid:      true,
		URL:        fmt.Sprintf("https://open.spotify.com/%s/%s", expectedType, input),
		ID:         input,
		EntityType: expectedType,
	}
}

// getEntityMismatchError はエンティティ種別不一致時のエラーメッセージを返します
func getEntityMismatchError(expectedType EntityType) string {
	switch expectedType {
	case EntityTrack:
		return "❌ Spotify の TrackURL を入力してください"
	case EntityArtist:
		return "❌ Spotify の ArtistURL を入力してください"
	case EntityAlbum:
		return "❌ Spotify の AlbumURL を入力してください"
	default:
		return "❌ Spotify の URL / ID として認識できませんでした。"
	}
}
