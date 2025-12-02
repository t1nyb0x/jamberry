package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/t1nyb0x/jamberry/internal/bot"
	"github.com/t1nyb0x/jamberry/internal/config"
	"github.com/t1nyb0x/jamberry/internal/handler"
	"github.com/t1nyb0x/jamberry/internal/infrastructure/cache"
	"github.com/t1nyb0x/jamberry/internal/infrastructure/tracktaste"
	"github.com/t1nyb0x/jamberry/internal/logger"
	"github.com/t1nyb0x/jamberry/internal/ratelimit"
	"github.com/t1nyb0x/jamberry/internal/usecase"
)

func main() {
	// 設定の読み込み
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// ロガーのセットアップ
	logger.Setup(cfg.LogLevel)

	slog.Info("starting jamberry", "log_level", cfg.LogLevel)

	// Botの作成
	b, err := bot.New(cfg.DiscordBotToken)
	if err != nil {
		slog.Error("failed to create bot", "error", err)
		os.Exit(1)
	}

	// インフラ層の作成
	ttClient := tracktaste.NewClient(cfg.TrackTasteAPIURL)
	cacheManager := cache.NewManager(cfg.RedisURL)
	defer cacheManager.Close()

	limiter := ratelimit.NewLimiter()

	// ユースケース層の作成
	trackUC := usecase.NewTrackUseCase(ttClient)
	artistUC := usecase.NewArtistUseCase(ttClient)
	albumUC := usecase.NewAlbumUseCase(ttClient)
	recommendUC := usecase.NewRecommendUseCase(ttClient)
	searchUC := usecase.NewSearchUseCase(ttClient)

	// ハンドラーの作成
	h := handler.NewHandler(
		trackUC,
		artistUC,
		albumUC,
		recommendUC,
		searchUC,
		cacheManager,
		limiter,
		ttClient,
	)

	// インタラクションハンドラーの登録
	b.AddHandler(h.HandleInteraction)

	// Botの起動
	if err := b.Start(); err != nil {
		slog.Error("failed to start bot", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := b.Stop(); err != nil {
			slog.Warn("failed to stop bot", "error", err)
		}
	}()

	// バックグラウンドタスクの開始
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// L1キャッシュのクリーンアップ
	cacheManager.StartL1Cleanup(ctx, 5*time.Minute)

	// レートリミッターのクリーンアップ
	done := make(chan struct{})
	limiter.StartCleanup(done, 30*time.Second)

	slog.Info("jamberry is now running. Press CTRL+C to exit.")

	// シグナル待機（Graceful Shutdown）
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	slog.Info("shutting down...")

	// クリーンアップ
	cancel()
	close(done)

	slog.Info("jamberry has been shut down")
}
