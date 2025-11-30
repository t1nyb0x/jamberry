package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/t1nyb0x/jamberry/internal/domain"
)

const (
	// L1TTL はL1キャッシュ（インメモリ）のTTLです
	L1TTL = 10 * time.Minute
	// L2TTL はL2キャッシュ（Redis）のTTLです
	L2TTL = 30 * 24 * time.Hour // 30日
)

// l1Entry はL1キャッシュのエントリを表します
type l1Entry struct {
	data      *domain.PaginationData
	expiresAt time.Time
}

// Manager はキャッシュマネージャーです
// domain.CacheRepository インターフェースを実装します
type Manager struct {
	l1       sync.Map
	redis    *redis.Client
	redisOK  bool
	redisMux sync.RWMutex
}

// インターフェース実装の確認
var _ domain.CacheRepository = (*Manager)(nil)

// NewManager は新しいキャッシュマネージャーを作成します
func NewManager(redisURL string) *Manager {
	m := &Manager{
		redisOK: false,
	}

	// Redis接続
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		slog.Warn("failed to parse redis URL, running without L2 cache", "error", err)
		return m
	}

	client := redis.NewClient(opts)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		slog.Warn("failed to connect to redis, running without L2 cache", "error", err)
		return m
	}

	m.redis = client
	m.redisOK = true
	slog.Info("connected to redis")

	return m
}

// makeKey はキャッシュキーを生成します
func makeKey(key string) string {
	return fmt.Sprintf("pagination:%s", key)
}

// Set はキャッシュにデータを保存します
func (m *Manager) Set(ctx context.Context, key string, data *domain.PaginationData) error {
	cacheKey := makeKey(key)

	// L1に保存
	m.l1.Store(cacheKey, &l1Entry{
		data:      data,
		expiresAt: time.Now().Add(L1TTL),
	})
	slog.Debug("cache set to L1", "key", cacheKey, "command", data.Command, "total", data.Total)

	// L2に保存
	if m.isRedisOK() {
		jsonData, err := json.Marshal(data)
		if err != nil {
			slog.Warn("failed to marshal cache data", "key", cacheKey, "error", err)
			return nil // L1には保存済みなのでエラーは返さない
		}

		if err := m.redis.Set(ctx, cacheKey, jsonData, L2TTL).Err(); err != nil {
			slog.Warn("failed to set cache in redis", "key", cacheKey, "error", err)
		} else {
			slog.Debug("cache set to L2", "key", cacheKey)
		}
	}

	return nil
}

// Get はキャッシュからデータを取得します
func (m *Manager) Get(ctx context.Context, key string) (*domain.PaginationData, error) {
	cacheKey := makeKey(key)

	// L1から取得
	if entry, ok := m.l1.Load(cacheKey); ok {
		l1e := entry.(*l1Entry)
		if time.Now().Before(l1e.expiresAt) {
			slog.Debug("cache hit L1", "key", cacheKey)
			return l1e.data, nil
		}
		// 期限切れの場合は削除
		m.l1.Delete(cacheKey)
		slog.Debug("cache expired L1", "key", cacheKey)
	}

	// L2から取得
	if m.isRedisOK() {
		jsonData, err := m.redis.Get(ctx, cacheKey).Bytes()
		if err == nil {
			var data domain.PaginationData
			if err := json.Unmarshal(jsonData, &data); err == nil {
				// L1に書き戻し
				m.l1.Store(cacheKey, &l1Entry{
					data:      &data,
					expiresAt: time.Now().Add(L1TTL),
				})
				slog.Debug("cache hit L2, restored to L1", "key", cacheKey)
				return &data, nil
			}
		} else if err != redis.Nil {
			slog.Warn("failed to get cache from redis", "key", cacheKey, "error", err)
		}
	}

	slog.Debug("cache miss", "key", cacheKey)
	return nil, fmt.Errorf("cache not found")
}

// Delete はキャッシュからデータを削除します
func (m *Manager) Delete(ctx context.Context, key string) {
	cacheKey := makeKey(key)
	m.l1.Delete(cacheKey)

	if m.isRedisOK() {
		if err := m.redis.Del(ctx, cacheKey).Err(); err != nil {
			slog.Warn("failed to delete cache from redis", "error", err)
		}
	}
}

// Close はキャッシュマネージャーをクローズします
func (m *Manager) Close() error {
	if m.redis != nil {
		return m.redis.Close()
	}
	return nil
}

// isRedisOK はRedisが利用可能かどうかを返します
func (m *Manager) isRedisOK() bool {
	m.redisMux.RLock()
	defer m.redisMux.RUnlock()
	return m.redisOK
}

// CleanupL1 は期限切れのL1キャッシュエントリを削除します
func (m *Manager) CleanupL1() {
	now := time.Now()
	m.l1.Range(func(key, value interface{}) bool {
		entry := value.(*l1Entry)
		if now.After(entry.expiresAt) {
			m.l1.Delete(key)
		}
		return true
	})
}

// StartL1Cleanup はL1キャッシュの定期クリーンアップを開始します
func (m *Manager) StartL1Cleanup(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				m.CleanupL1()
			}
		}
	}()
}
