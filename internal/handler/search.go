package handler

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/t1nyb0x/jamberry/internal/domain"
	"github.com/t1nyb0x/jamberry/internal/presenter"
	"github.com/t1nyb0x/jamberry/internal/usecase"
)

// handleSearch は検索コマンドを処理します
func (h *Handler) handleSearch(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) {
	if len(options) == 0 {
		slog.Info("validation failed: empty input", "command", "jam search")
		h.responder.RespondEphemeral(s, i, "❌ 検索キーワードを入力してください。")
		return
	}

	query := options[0].StringValue()

	// DeferReply
	if err := h.responder.DeferReply(s, i); err != nil {
		slog.Error("failed to defer reply", "command", "jam search", "error", err)
		return
	}

	ctx := context.Background()
	output, err := h.searchUseCase.SearchTracks(ctx, usecase.SearchInput{Query: query})
	if err != nil {
		h.responder.EditResponse(s, i, err.Error())
		return
	}

	userID := getUserID(i)
	totalPages := (len(output.Tracks) + PageSize - 1) / PageSize
	emb := presenter.BuildSearchEmbed(output.Query, output.Tracks, 0, PageSize, len(output.Tracks))

	// 初期ボタン（placeholderで仮設定）
	components := []discordgo.MessageComponent{
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
	}

	msg, err := h.responder.EditResponseWithComponents(s, i, emb, components)
	if err != nil {
		slog.Error("failed to send response", "error", err)
		return
	}

	// キャッシュに保存
	itemsJSON, _ := json.Marshal(output.Tracks)
	cacheData := &domain.PaginationData{
		Command: "search",
		Query:   output.Query,
		Type:    "track",
		Items:   itemsJSON,
		Total:   len(output.Tracks),
		OwnerID: userID,
	}
	if err := h.cache.Set(ctx, msg.ID, cacheData); err != nil {
		slog.Warn("failed to cache data", "error", err)
	}

	// ボタンのCustomIDを更新
	updatedComponents := presenter.BuildPaginationButtons(msg.ID, 0, totalPages)
	_, _ = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Components: &updatedComponents,
	})

	slog.Info("command completed", "command", "search", "query", output.Query, "result_count", len(output.Tracks), "message_id", msg.ID)
}
