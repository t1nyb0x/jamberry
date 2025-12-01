package domain

import (
	"context"
	"encoding/json"
)

// PaginationData はページネーション用のキャッシュデータを表します
type PaginationData struct {
	Command string          `json:"command"`
	Query   string          `json:"query"`
	Type    string          `json:"type"`
	Items   json.RawMessage `json:"items"`
	Total   int             `json:"total"`
	OwnerID string          `json:"owner_id"`
	Mode    string          `json:"mode,omitempty"` // レコメンドモード（recommend専用）
}

// CacheRepository はキャッシュを操作するリポジトリインターフェースです
type CacheRepository interface {
	// Set はキャッシュにデータを保存します
	Set(ctx context.Context, key string, data *PaginationData) error

	// Get はキャッシュからデータを取得します
	Get(ctx context.Context, key string) (*PaginationData, error)

	// Delete はキャッシュからデータを削除します
	Delete(ctx context.Context, key string)
}
