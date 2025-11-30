package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/t1nyb0x/jamberry/internal/domain"
	"github.com/t1nyb0x/jamberry/internal/presenter"
)

// handlePaging はページングボタンを処理します
func (h *Handler) handlePaging(s *discordgo.Session, i *discordgo.InteractionCreate, messageID, action string, parts []string, userID string) {
	ctx := context.Background()

	// キャッシュからデータを取得
	cacheData, err := h.cache.Get(ctx, messageID)
	if err != nil {
		slog.Info("cache expired for button interaction", "message_id", messageID, "user_id", userID)
		h.responder.RespondEphemeral(s, i, "データの有効期限が切れました。再度コマンドを実行してください。")
		return
	}

	// 操作権限チェック
	if cacheData.OwnerID != userID {
		slog.Info("paging permission denied", "action", action, "owner_id", cacheData.OwnerID, "user_id", userID)
		h.responder.RespondEphemeral(s, i, "この操作はコマンド実行者のみが使用できます。『👁 自分も見る』ボタンを押すと、あなた専用の表示ができます。")
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
	emb := buildEmbedFromCache(cacheData, newPage)
	components := presenter.BuildPaginationButtons(messageID, newPage, totalPages)

	// メッセージを更新
	h.responder.UpdateMessage(s, i, emb, components)
	slog.Debug("page updated", "action", action, "message_id", messageID, "page", newPage, "total_pages", totalPages)
}

// handleViewOwn は「自分も見る」ボタンを処理します
func (h *Handler) handleViewOwn(s *discordgo.Session, i *discordgo.InteractionCreate, messageID string) {
	ctx := context.Background()
	userID := getUserID(i)

	slog.Debug("view_own button pressed", "message_id", messageID, "user_id", userID)

	// キャッシュからデータを取得
	cacheData, err := h.cache.Get(ctx, messageID)
	if err != nil {
		slog.Info("cache expired for view_own", "message_id", messageID, "user_id", userID)
		h.responder.RespondEphemeral(s, i, "データの有効期限が切れました。再度コマンドを実行してください。")
		return
	}

	totalPages := (cacheData.Total + PageSize - 1) / PageSize
	emb := buildEmbedFromCache(cacheData, 0)

	components := []discordgo.MessageComponent{
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
	}

	h.responder.RespondEphemeralWithEmbed(s, i, emb, components)
}

// buildEmbedFromCache はキャッシュデータからEmbedを構築します
func buildEmbedFromCache(cacheData *domain.PaginationData, page int) *discordgo.MessageEmbed {
	if cacheData.Command == "recommend" {
		var items []domain.SimilarTrack
		json.Unmarshal(cacheData.Items, &items)
		return presenter.BuildRecommendEmbed(cacheData.Query, items, page, PageSize, cacheData.Total)
	}

	var items []domain.Track
	json.Unmarshal(cacheData.Items, &items)
	return presenter.BuildSearchEmbed(cacheData.Query, items, page, PageSize, cacheData.Total)
}
