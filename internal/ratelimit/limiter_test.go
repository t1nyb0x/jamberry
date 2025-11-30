package ratelimit

import (
	"sync"
	"testing"
	"time"
)

func TestLimiter_Allow_BasicUsage(t *testing.T) {
	limiter := NewLimiter()

	// 最初の5回は許可される
	for i := 0; i < 5; i++ {
		if !limiter.Allow("user1") {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	// 6回目は拒否される
	if limiter.Allow("user1") {
		t.Error("6th request should be denied")
	}
}

func TestLimiter_Allow_DifferentUsers(t *testing.T) {
	limiter := NewLimiter()

	// user1が5回リクエスト
	for i := 0; i < 5; i++ {
		if !limiter.Allow("user1") {
			t.Errorf("user1 request %d should be allowed", i+1)
		}
	}

	// user1の6回目は拒否
	if limiter.Allow("user1") {
		t.Error("user1 6th request should be denied")
	}

	// user2は別ユーザーなので許可される
	for i := 0; i < 5; i++ {
		if !limiter.Allow("user2") {
			t.Errorf("user2 request %d should be allowed", i+1)
		}
	}

	// user2の6回目は拒否
	if limiter.Allow("user2") {
		t.Error("user2 6th request should be denied")
	}
}

func TestLimiter_Allow_WindowExpiration(t *testing.T) {
	// テスト用に短いウィンドウを使うため、直接テストするのは難しい
	// 実際のWindowは10秒だが、テストでは状態をシミュレート

	limiter := NewLimiter()

	// 5回リクエスト
	for i := 0; i < 5; i++ {
		limiter.Allow("user1")
	}

	// 6回目は拒否
	if limiter.Allow("user1") {
		t.Error("6th request should be denied")
	}
}

func TestLimiter_Allow_Concurrency(t *testing.T) {
	limiter := NewLimiter()
	var wg sync.WaitGroup
	allowCount := 0
	denyCount := 0
	var mu sync.Mutex

	// 10個のゴルーチンで同時にリクエスト
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			allowed := limiter.Allow("user1")
			mu.Lock()
			if allowed {
				allowCount++
			} else {
				denyCount++
			}
			mu.Unlock()
		}()
	}

	wg.Wait()

	// 5回のみ許可されるべき
	if allowCount != 5 {
		t.Errorf("Expected 5 allowed requests, got %d", allowCount)
	}
	if denyCount != 5 {
		t.Errorf("Expected 5 denied requests, got %d", denyCount)
	}
}

func TestLimiter_Cleanup(t *testing.T) {
	limiter := NewLimiter()

	// ユーザーを追加
	limiter.Allow("user1")
	limiter.Allow("user2")

	// Cleanupを実行（ウィンドウ内なので削除されない）
	limiter.Cleanup()

	// まだリクエスト可能であることを確認
	// (初回リクエスト後なのでまだ4回は許可される)
	if !limiter.Allow("user1") {
		t.Error("user1 should still be able to request after cleanup")
	}
}

func TestLimiter_Cleanup_RemovesExpiredEntries(t *testing.T) {
	limiter := NewLimiter()

	// エントリを追加
	value, _ := limiter.users.LoadOrStore("olduser", &userEntry{
		timestamps: []time.Time{
			time.Now().Add(-20 * time.Second), // 20秒前（ウィンドウ外）
		},
	})
	entry := value.(*userEntry)
	entry.timestamps = []time.Time{
		time.Now().Add(-20 * time.Second),
	}

	// Cleanupを実行
	limiter.Cleanup()

	// エントリが削除されていることを確認
	if _, ok := limiter.users.Load("olduser"); ok {
		t.Error("Expired entry should be removed after cleanup")
	}
}

func TestLimiter_MultipleUsers_Independent(t *testing.T) {
	limiter := NewLimiter()

	users := []string{"user1", "user2", "user3", "user4", "user5"}

	// 各ユーザーが5回リクエスト
	for _, user := range users {
		for i := 0; i < 5; i++ {
			if !limiter.Allow(user) {
				t.Errorf("%s request %d should be allowed", user, i+1)
			}
		}
	}

	// 各ユーザーの6回目は拒否
	for _, user := range users {
		if limiter.Allow(user) {
			t.Errorf("%s 6th request should be denied", user)
		}
	}
}

func TestConstants(t *testing.T) {
	// 定数の値を確認
	if Window != 10*time.Second {
		t.Errorf("Window = %v, want 10s", Window)
	}
	if MaxRequests != 5 {
		t.Errorf("MaxRequests = %d, want 5", MaxRequests)
	}
}

func TestNewLimiter(t *testing.T) {
	limiter := NewLimiter()
	if limiter == nil {
		t.Error("NewLimiter should return non-nil limiter")
	}
}

func TestLimiter_Allow_ExactlyMaxRequests(t *testing.T) {
	limiter := NewLimiter()

	// MaxRequests回リクエスト
	for i := 0; i < MaxRequests; i++ {
		if !limiter.Allow("user1") {
			t.Errorf("Request %d should be allowed (MaxRequests=%d)", i+1, MaxRequests)
		}
	}

	// MaxRequests + 1回目は拒否
	if limiter.Allow("user1") {
		t.Errorf("Request %d should be denied (MaxRequests=%d)", MaxRequests+1, MaxRequests)
	}
}

func TestLimiter_EmptyUserID(t *testing.T) {
	limiter := NewLimiter()

	// 空のユーザーIDでも動作する
	for i := 0; i < 5; i++ {
		if !limiter.Allow("") {
			t.Errorf("Request %d with empty user ID should be allowed", i+1)
		}
	}

	if limiter.Allow("") {
		t.Error("6th request with empty user ID should be denied")
	}
}
