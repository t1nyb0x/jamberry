package ratelimit

import (
	"sync"
	"time"
)

const (
	// Window はレートリミットのウィンドウサイズです
	Window = 10 * time.Second
	// MaxRequests はウィンドウ内の最大リクエスト数です
	MaxRequests = 5
)

// Limiter はユーザーごとのレートリミッターです
type Limiter struct {
	users sync.Map
}

// userEntry はユーザーごとのレートリミット情報を保持します
type userEntry struct {
	mu         sync.Mutex
	timestamps []time.Time
}

// NewLimiter は新しいレートリミッターを作成します
func NewLimiter() *Limiter {
	return &Limiter{}
}

// Allow は指定されたユーザーのリクエストを許可するかどうかを判定します
func (l *Limiter) Allow(userID string) bool {
	now := time.Now()
	windowStart := now.Add(-Window)

	// ユーザーエントリを取得または作成
	value, _ := l.users.LoadOrStore(userID, &userEntry{})
	entry := value.(*userEntry)

	entry.mu.Lock()
	defer entry.mu.Unlock()

	// ウィンドウ外のタイムスタンプを削除
	var validTimestamps []time.Time
	for _, ts := range entry.timestamps {
		if ts.After(windowStart) {
			validTimestamps = append(validTimestamps, ts)
		}
	}

	// 制限を超えていないかチェック
	if len(validTimestamps) >= MaxRequests {
		entry.timestamps = validTimestamps
		return false
	}

	// タイムスタンプを追加
	entry.timestamps = append(validTimestamps, now)
	return true
}

// Cleanup は古いエントリをクリーンアップします
func (l *Limiter) Cleanup() {
	now := time.Now()
	windowStart := now.Add(-Window)

	l.users.Range(func(key, value interface{}) bool {
		entry := value.(*userEntry)
		entry.mu.Lock()

		var validTimestamps []time.Time
		for _, ts := range entry.timestamps {
			if ts.After(windowStart) {
				validTimestamps = append(validTimestamps, ts)
			}
		}

		if len(validTimestamps) == 0 {
			entry.mu.Unlock()
			l.users.Delete(key)
		} else {
			entry.timestamps = validTimestamps
			entry.mu.Unlock()
		}

		return true
	})
}

// StartCleanup は定期クリーンアップを開始します
func (l *Limiter) StartCleanup(done <-chan struct{}, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-done:
				ticker.Stop()
				return
			case <-ticker.C:
				l.Cleanup()
			}
		}
	}()
}
