package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/t1nyb0x/jamberry/internal/cache"
	"github.com/t1nyb0x/jamberry/internal/embed"
	"github.com/t1nyb0x/jamberry/internal/ratelimit"
	"github.com/t1nyb0x/jamberry/internal/spotify"
	"github.com/t1nyb0x/jamberry/internal/tracktaste"
)

const (
	// PageSize は1ページあたりの表示件数です
	PageSize = 5
)

// Handler はDiscordコマンドハンドラーです
type Handler struct {
	ttClient  *tracktaste.Client
	cache     *cache.Manager
	limiter   *ratelimit.Limiter
}

// NewHandler は新しいハンドラーを作成します
func NewHandler(ttClient *tracktaste.Client, cache *cache.Manager, limiter *ratelimit.Limiter) *Handler {
	return &Handler{
		ttClient:  ttClient,
		cache:     cache,
		limiter:   limiter,
	}
}

// Commands はスラッシュコマンドの定義を返します
func Commands() []*discordgo.ApplicationCommand {
	urlOption := &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "url",
		Description: "Spotify の URL, URI, または ID を入力",
		Required:    true,
	}

	return []*discordgo.ApplicationCommand{
		{
			Name:        "track",
			Description: "Spotifyトラックの詳細情報を取得します",
			Options: []*discordgo.ApplicationCommandOption{
				urlOption,
			},
		},
		{
			Name:        "artist",
			Description: "Spotifyアーティストの詳細情報を取得します",
			Options: []*discordgo.ApplicationCommandOption{
				urlOption,
			},
		},
		{
			Name:        "album",
			Description: "Spotifyアルバムの詳細情報を取得します",
			Options: []*discordgo.ApplicationCommandOption{
				urlOption,
			},
		},
		{
			Name:        "recommend",
			Description: "トラックに基づくおすすめ楽曲を取得します",
			Options: []*discordgo.ApplicationCommandOption{
				urlOption,
			},
		},
		{
			Name:        "search",
			Description: "Spotifyでトラックを検索します",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "query",
					Description: "検索キーワードを入力",
					Required:    true,
				},
			},
		},
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
	// ユーザーIDを取得
	var userID string
	if i.Member != nil {
		userID = i.Member.User.ID
	} else if i.User != nil {
		userID = i.User.ID
	}

	// ログ出力
	var guildID, channelID string
	if i.GuildID != "" {
		guildID = i.GuildID
	}
	if i.ChannelID != "" {
		channelID = i.ChannelID
	}
	cmdName := i.ApplicationCommandData().Name
	slog.Info("command received",
		"guild_id", guildID,
		"channel_id", channelID,
		"command", cmdName,
		"user_id", userID,
	)

	// レートリミットチェック
	if !h.limiter.Allow(userID) {
		slog.Warn("rate limit exceeded", "user_id", userID)
		h.respondEphemeral(s, i, "⏳ 少し待ってから再試行してください。")
		return
	}

	// コマンドごとの処理
	switch cmdName {
	case "track":
		h.handleTrack(s, i)
	case "artist":
		h.handleArtist(s, i)
	case "album":
		h.handleAlbum(s, i)
	case "recommend":
		h.handleRecommend(s, i)
	case "search":
		h.handleSearch(s, i)
	}
}

// handleTrack はトラック情報取得コマンドを処理します
func (h *Handler) handleTrack(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		slog.Info("validation failed: empty input", "command", "track")
		h.respondEphemeral(s, i, "❌ URL を入力してください。")
		return
	}

	input := options[0].StringValue()

	// バリデーション
	result := spotify.ValidateInput(input, spotify.EntityTrack)
	if !result.Valid {
		slog.Info("validation failed", "command", "track", "input", input, "error", result.Error)
		h.respondEphemeral(s, i, result.Error)
		return
	}
	slog.Debug("validation passed", "command", "track", "url", result.URL, "id", result.ID)

	// DeferReply
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	}); err != nil {
		slog.Error("failed to defer reply", "command", "track", "error", err)
		return
	}

	// tracktaste API呼び出し
	ctx := context.Background()
	track, err := h.ttClient.FetchTrack(ctx, result.URL)
	if err != nil {
		slog.Warn("track fetch failed", "command", "track", "url", result.URL, "error", err)
		h.editResponse(s, i, err.Error())
		return
	}

	// Embed構築・返信
	emb := embed.BuildTrackEmbed(track)
	h.editResponseEmbed(s, i, emb)
	slog.Info("command completed", "command", "track", "track_name", track.Name, "track_id", track.ID)
}

// handleArtist はアーティスト情報取得コマンドを処理します
func (h *Handler) handleArtist(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		slog.Info("validation failed: empty input", "command", "artist")
		h.respondEphemeral(s, i, "❌ URL を入力してください。")
		return
	}

	input := options[0].StringValue()

	// バリデーション
	result := spotify.ValidateInput(input, spotify.EntityArtist)
	if !result.Valid {
		slog.Info("validation failed", "command", "artist", "input", input, "error", result.Error)
		h.respondEphemeral(s, i, result.Error)
		return
	}
	slog.Debug("validation passed", "command", "artist", "url", result.URL, "id", result.ID)

	// DeferReply
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	}); err != nil {
		slog.Error("failed to defer reply", "command", "artist", "error", err)
		return
	}

	// tracktaste API呼び出し
	ctx := context.Background()
	artist, err := h.ttClient.FetchArtist(ctx, result.URL)
	if err != nil {
		slog.Warn("artist fetch failed", "command", "artist", "url", result.URL, "error", err)
		h.editResponse(s, i, err.Error())
		return
	}

	// Embed構築・返信
	emb := embed.BuildArtistEmbed(artist)
	h.editResponseEmbed(s, i, emb)
	slog.Info("command completed", "command", "artist", "artist_name", artist.Name, "artist_id", artist.ID)
}

// handleAlbum はアルバム情報取得コマンドを処理します
func (h *Handler) handleAlbum(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		slog.Info("validation failed: empty input", "command", "album")
		h.respondEphemeral(s, i, "❌ URL を入力してください。")
		return
	}

	input := options[0].StringValue()

	// バリデーション
	result := spotify.ValidateInput(input, spotify.EntityAlbum)
	if !result.Valid {
		slog.Info("validation failed", "command", "album", "input", input, "error", result.Error)
		h.respondEphemeral(s, i, result.Error)
		return
	}
	slog.Debug("validation passed", "command", "album", "url", result.URL, "id", result.ID)

	// DeferReply
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	}); err != nil {
		slog.Error("failed to defer reply", "command", "album", "error", err)
		return
	}

	// tracktaste API呼び出し
	ctx := context.Background()
	album, err := h.ttClient.FetchAlbum(ctx, result.URL)
	if err != nil {
		slog.Warn("album fetch failed", "command", "album", "url", result.URL, "error", err)
		h.editResponse(s, i, err.Error())
		return
	}

	// Embed構築・返信
	emb := embed.BuildAlbumEmbed(album)
	h.editResponseEmbed(s, i, emb)
	slog.Info("command completed", "command", "album", "album_name", album.Name, "album_id", album.ID)
}

// handleRecommend はレコメンド取得コマンドを処理します
func (h *Handler) handleRecommend(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		slog.Info("validation failed: empty input", "command", "recommend")
		h.respondEphemeral(s, i, "❌ URL を入力してください。")
		return
	}

	input := options[0].StringValue()

	// バリデーション（トラックURLのみ受け付ける）
	result := spotify.ValidateInput(input, spotify.EntityTrack)
	if !result.Valid {
		slog.Info("validation failed", "command", "recommend", "input", input, "error", result.Error)
		h.respondEphemeral(s, i, result.Error)
		return
	}
	slog.Debug("validation passed", "command", "recommend", "url", result.URL, "id", result.ID)

	// DeferReply
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	}); err != nil {
		slog.Error("failed to defer reply", "command", "recommend", "error", err)
		return
	}

	ctx := context.Background()

	// 元のトラック情報を取得（表示用）
	track, err := h.ttClient.FetchTrack(ctx, result.URL)
	if err != nil {
		slog.Warn("track fetch failed for recommend", "command", "recommend", "url", result.URL, "error", err)
		h.editResponse(s, i, err.Error())
		return
	}

	// 類似トラックを取得
	similar, err := h.ttClient.FetchSimilar(ctx, result.URL)
	if err != nil {
		slog.Warn("similar fetch failed", "command", "recommend", "url", result.URL, "error", err)
		h.editResponse(s, i, err.Error())
		return
	}

	if len(similar) == 0 {
		slog.Info("no results found", "command", "recommend", "track_name", track.Name)
		h.editResponse(s, i, "🔍 該当する結果は見つかりませんでした。")
		return
	}

	// ユーザーIDを取得
	var userID string
	if i.Member != nil {
		userID = i.Member.User.ID
	} else if i.User != nil {
		userID = i.User.ID
	}

	// レスポンスを送信してメッセージIDを取得
	totalPages := (len(similar) + PageSize - 1) / PageSize
	emb := embed.BuildRecommendEmbed(track.Name, similar, 0, PageSize, len(similar))

	msg, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{emb},
		Components: &[]discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "◀ 前へ",
						Style:    discordgo.SecondaryButton,
						CustomID: "page_prev:placeholder:0",
						Disabled: true,
					},
					discordgo.Button{
						Label:    "次へ ▶",
						Style:    discordgo.SecondaryButton,
						CustomID: "page_next:placeholder:0",
						Disabled: totalPages <= 1,
					},
					discordgo.Button{
						Label:    "👁 自分も見る",
						Style:    discordgo.PrimaryButton,
						CustomID: "view_own:placeholder",
					},
				},
			},
		},
	})

	if err != nil {
		slog.Error("failed to send response", "error", err)
		return
	}

	// キャッシュに保存
	itemsJSON, _ := json.Marshal(similar)
	cacheData := &cache.CacheData{
		Command:   "recommend",
		Query:     track.Name,
		Type:      "track",
		Items:     json.RawMessage(itemsJSON),
		Total:     len(similar),
		OwnerID:   userID,
	}
	if err := h.cache.Set(ctx, msg.ID, cacheData); err != nil {
		slog.Warn("failed to cache data", "error", err)
	}

	// ボタンのCustomIDを更新
	components := embed.BuildPaginationButtons(msg.ID, 0, totalPages)
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Components: &components,
	})
	slog.Info("command completed", "command", "recommend", "track_name", track.Name, "result_count", len(similar), "message_id", msg.ID)
}

// handleSearch は検索コマンドを処理します
func (h *Handler) handleSearch(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		slog.Info("validation failed: empty input", "command", "search")
		h.respondEphemeral(s, i, "❌ 検索キーワードを入力してください。")
		return
	}

	query := strings.TrimSpace(options[0].StringValue())
	if query == "" {
		slog.Info("validation failed: empty query", "command", "search")
		h.respondEphemeral(s, i, "❌ 検索キーワードを入力してください。")
		return
	}
	slog.Debug("search query received", "command", "search", "query", query)

	// DeferReply
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	}); err != nil {
		slog.Error("failed to defer reply", "command", "search", "error", err)
		return
	}

	ctx := context.Background()

	// 検索実行
	results, err := h.ttClient.SearchTracks(ctx, query)
	if err != nil {
		slog.Warn("search failed", "command", "search", "query", query, "error", err)
		h.editResponse(s, i, err.Error())
		return
	}

	if len(results) == 0 {
		slog.Info("no results found", "command", "search", "query", query)
		h.editResponse(s, i, "🔍 該当する結果は見つかりませんでした。")
		return
	}

	// ユーザーIDを取得
	var userID string
	if i.Member != nil {
		userID = i.Member.User.ID
	} else if i.User != nil {
		userID = i.User.ID
	}

	// レスポンスを送信
	totalPages := (len(results) + PageSize - 1) / PageSize
	emb := embed.BuildSearchEmbed(query, results, 0, PageSize, len(results))

	msg, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{emb},
		Components: &[]discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "◀ 前へ",
						Style:    discordgo.SecondaryButton,
						CustomID: "page_prev:placeholder:0",
						Disabled: true,
					},
					discordgo.Button{
						Label:    "次へ ▶",
						Style:    discordgo.SecondaryButton,
						CustomID: "page_next:placeholder:0",
						Disabled: totalPages <= 1,
					},
					discordgo.Button{
						Label:    "👁 自分も見る",
						Style:    discordgo.PrimaryButton,
						CustomID: "view_own:placeholder",
					},
				},
			},
		},
	})

	if err != nil {
		slog.Error("failed to send response", "error", err)
		return
	}

	// キャッシュに保存
	itemsJSON, _ := json.Marshal(results)
	cacheData := &cache.CacheData{
		Command:   "search",
		Query:     query,
		Type:      "track",
		Items:     json.RawMessage(itemsJSON),
		Total:     len(results),
		OwnerID:   userID,
	}
	if err := h.cache.Set(ctx, msg.ID, cacheData); err != nil {
		slog.Warn("failed to cache data", "error", err)
	}

	// ボタンのCustomIDを更新
	components := embed.BuildPaginationButtons(msg.ID, 0, totalPages)
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Components: &components,
	})
	slog.Info("command completed", "command", "search", "query", query, "result_count", len(results), "message_id", msg.ID)
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

	// ユーザーIDを取得
	var userID string
	if i.Member != nil {
		userID = i.Member.User.ID
	} else if i.User != nil {
		userID = i.User.ID
	}

	slog.Debug("button interaction received", "action", action, "message_id", messageID, "user_id", userID)

	ctx := context.Background()

	// キャッシュからデータを取得
	cacheData, err := h.cache.Get(ctx, messageID)
	if err != nil {
		slog.Info("cache expired for button interaction", "message_id", messageID, "user_id", userID)
		h.respondEphemeral(s, i, "データの有効期限が切れました。再度コマンドを実行してください。")
		return
	}

	switch action {
	case "page_prev", "page_next":
		h.handlePaging(s, i, cacheData, action, parts, userID)
	case "view_own":
		h.handleViewOwn(s, i, cacheData, messageID)
	}
}

// handlePaging はページングボタンを処理します
func (h *Handler) handlePaging(s *discordgo.Session, i *discordgo.InteractionCreate, cacheData *cache.CacheData, action string, parts []string, userID string) {
	// 操作権限チェック
	if cacheData.OwnerID != userID {
		slog.Info("paging permission denied", "action", action, "owner_id", cacheData.OwnerID, "user_id", userID)
		h.respondEphemeral(s, i, "この操作はコマンド実行者のみが使用できます。『👁 自分も見る』ボタンを押すと、あなた専用の表示ができます。")
		return
	}

	// 現在のページを取得
	currentPage := 0
	if len(parts) >= 3 {
		if p, err := strconv.Atoi(parts[2]); err == nil {
			currentPage = p
		}
	}

	// 新しいページを計算
	newPage := currentPage
	if action == "page_prev" && currentPage > 0 {
		newPage = currentPage - 1
	} else if action == "page_next" {
		newPage = currentPage + 1
	}

	totalPages := (cacheData.Total + PageSize - 1) / PageSize
	if newPage >= totalPages {
		newPage = totalPages - 1
	}
	if newPage < 0 {
		newPage = 0
	}

	// Embedを構築
	var emb *discordgo.MessageEmbed
	if cacheData.Command == "recommend" {
		var items []tracktaste.SimilarTrack
		json.Unmarshal(cacheData.Items.(json.RawMessage), &items)
		emb = embed.BuildRecommendEmbed(cacheData.Query, items, newPage, PageSize, cacheData.Total)
	} else {
		var items []tracktaste.SearchTrack
		json.Unmarshal(cacheData.Items.(json.RawMessage), &items)
		emb = embed.BuildSearchEmbed(cacheData.Query, items, newPage, PageSize, cacheData.Total)
	}

	messageID := parts[1]
	components := embed.BuildPaginationButtons(messageID, newPage, totalPages)

	// メッセージを更新
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{emb},
			Components: components,
		},
	})
	slog.Debug("page updated", "action", action, "message_id", messageID, "page", newPage, "total_pages", totalPages)
}

// handleViewOwn は「自分も見る」ボタンを処理します
func (h *Handler) handleViewOwn(s *discordgo.Session, i *discordgo.InteractionCreate, cacheData *cache.CacheData, messageID string) {
	var userID string
	if i.Member != nil {
		userID = i.Member.User.ID
	} else if i.User != nil {
		userID = i.User.ID
	}
	slog.Debug("view_own button pressed", "message_id", messageID, "user_id", userID, "command", cacheData.Command)

	totalPages := (cacheData.Total + PageSize - 1) / PageSize

	// Embedを構築
	var emb *discordgo.MessageEmbed
	if cacheData.Command == "recommend" {
		var items []tracktaste.SimilarTrack
		json.Unmarshal(cacheData.Items.(json.RawMessage), &items)
		emb = embed.BuildRecommendEmbed(cacheData.Query, items, 0, PageSize, cacheData.Total)
	} else {
		var items []tracktaste.SearchTrack
		json.Unmarshal(cacheData.Items.(json.RawMessage), &items)
		emb = embed.BuildSearchEmbed(cacheData.Query, items, 0, PageSize, cacheData.Total)
	}

	// Ephemeralで応答
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{emb},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "◀ 前へ",
							Style:    discordgo.SecondaryButton,
							CustomID: fmt.Sprintf("ephemeral_prev:%s:0", messageID),
							Disabled: true,
						},
						discordgo.Button{
							Label:    "次へ ▶",
							Style:    discordgo.SecondaryButton,
							CustomID: fmt.Sprintf("ephemeral_next:%s:0", messageID),
							Disabled: totalPages <= 1,
						},
					},
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
}

// respondEphemeral はEphemeralメッセージで応答します
func (h *Handler) respondEphemeral(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

// editResponse はDeferred応答を編集します
func (h *Handler) editResponse(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
}

// editResponseEmbed はDeferred応答をEmbedで編集します
func (h *Handler) editResponseEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, emb *discordgo.MessageEmbed) {
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{emb},
	})
}
