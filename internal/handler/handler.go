package handler

import (
	"log/slog"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/t1nyb0x/jamberry/internal/domain"
	"github.com/t1nyb0x/jamberry/internal/ratelimit"
	"github.com/t1nyb0x/jamberry/internal/usecase"
)

// Handler はDiscordコマンドハンドラーです
type Handler struct {
	trackUseCase     *usecase.TrackUseCase
	artistUseCase    *usecase.ArtistUseCase
	albumUseCase     *usecase.AlbumUseCase
	recommendUseCase *usecase.RecommendUseCase
	searchUseCase    *usecase.SearchUseCase
	cache            domain.CacheRepository
	limiter          *ratelimit.Limiter
	responder        *Responder
}

// NewHandler は新しいハンドラーを作成します
func NewHandler(
	trackUC *usecase.TrackUseCase,
	artistUC *usecase.ArtistUseCase,
	albumUC *usecase.AlbumUseCase,
	recommendUC *usecase.RecommendUseCase,
	searchUC *usecase.SearchUseCase,
	cache domain.CacheRepository,
	limiter *ratelimit.Limiter,
) *Handler {
	return &Handler{
		trackUseCase:     trackUC,
		artistUseCase:    artistUC,
		albumUseCase:     albumUC,
		recommendUseCase: recommendUC,
		searchUseCase:    searchUC,
		cache:            cache,
		limiter:          limiter,
		responder:        NewResponder(),
	}
}

// HandleInteraction はインタラクションを処理します
func (h *Handler) HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		h.handleCommand(s, i)
	case discordgo.InteractionMessageComponent:
		h.handleComponent(s, i)
	}
}

// handleCommand はスラッシュコマンドを処理します
func (h *Handler) handleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userID := getUserID(i)

	// サブコマンド名を取得
	cmdData := i.ApplicationCommandData()
	cmdName := cmdData.Name

	var subCmdName string
	var options []*discordgo.ApplicationCommandInteractionDataOption

	// サブコマンドの場合
	if cmdName == "jam" && len(cmdData.Options) > 0 {
		subCmdName = cmdData.Options[0].Name
		options = cmdData.Options[0].Options
	}

	// ログ出力
	slog.Info("command received",
		"guild_id", i.GuildID,
		"channel_id", i.ChannelID,
		"command", cmdName,
		"subcommand", subCmdName,
		"user_id", userID,
	)

	// レートリミットチェック
	if !h.limiter.Allow(userID) {
		slog.Warn("rate limit exceeded", "user_id", userID)
		h.responder.RespondEphemeral(s, i, "⏳ 少し待ってから再試行してください。")
		return
	}

	// サブコマンドごとの処理
	switch subCmdName {
	case "track":
		h.handleTrack(s, i, options)
	case "artist":
		h.handleArtist(s, i, options)
	case "album":
		h.handleAlbum(s, i, options)
	case "recommend":
		h.handleRecommend(s, i, options)
	case "search":
		h.handleSearch(s, i, options)
	}
}

// handleComponent はボタンコンポーネントを処理します
func (h *Handler) handleComponent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	customID := i.MessageComponentData().CustomID
	parts := strings.Split(customID, ":")

	if len(parts) < 2 {
		slog.Warn("invalid component custom_id", "custom_id", customID)
		return
	}

	action := parts[0]
	messageID := parts[1]
	userID := getUserID(i)

	slog.Debug("button interaction received", "action", action, "message_id", messageID, "user_id", userID)

	switch action {
	case "page_prev", "page_next":
		h.handlePaging(s, i, messageID, action, parts, userID)
	case "view_own":
		h.handleViewOwn(s, i, messageID)
	case "ephemeral_prev", "ephemeral_next":
		h.handleEphemeralPaging(s, i, messageID, action, parts)
	}
}

// getUserID はインタラクションからユーザーIDを取得します
func getUserID(i *discordgo.InteractionCreate) string {
	if i.Member != nil {
		return i.Member.User.ID
	}
	if i.User != nil {
		return i.User.ID
	}
	return ""
}
