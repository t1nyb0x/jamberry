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

const (
	// PageSize は1ページあたりの表示件数です
	PageSize = 5
)

// handleRecommend はレコメンド取得コマンドを処理します
func (h *Handler) handleRecommend(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) {
	if len(options) == 0 {
		slog.Info("validation failed: empty input", "command", "jam recommend")
		h.responder.RespondEphemeral(s, i, "❌ URL を入力してください。")
		return
	}

	// オプションの解析
	var input string
	var mode domain.RecommendMode

	for _, opt := range options {
		switch opt.Name {
		case "url":
			input = opt.StringValue()
		case "mode":
			mode = domain.RecommendMode(opt.StringValue())
		}
	}

	if input == "" {
		slog.Info("validation failed: empty input", "command", "jam recommend")
		h.responder.RespondEphemeral(s, i, "❌ URL を入力してください。")
		return
	}

	// DeferReply
	if err := h.responder.DeferReply(s, i); err != nil {
		slog.Error("failed to defer reply", "command", "jam recommend", "error", err)
		return
	}

	ctx := context.Background()
	output, err := h.recommendUseCase.GetRecommend(ctx, usecase.RecommendInput{
		Input: input,
		Mode:  mode,
	})
	if err != nil {
		h.responder.EditResponse(s, i, err.Error())
		return
	}

	userID := getUserID(i)
	totalPages := (len(output.Items) + PageSize - 1) / PageSize
	emb := presenter.BuildRecommendEmbed(output.SeedTrack.Name, output.Items, 0, PageSize, len(output.Items), output.Mode)

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
	itemsJSON, _ := json.Marshal(output.Items)
	cacheData := &domain.PaginationData{
		Command: "recommend",
		Query:   output.SeedTrack.Name,
		Type:    "track",
		Items:   itemsJSON,
		Total:   len(output.Items),
		OwnerID: userID,
		Mode:    string(output.Mode),
	}
	if err := h.cache.Set(ctx, msg.ID, cacheData); err != nil {
		slog.Warn("failed to cache data", "error", err)
	}

	// ボタンのCustomIDを更新
	updatedComponents := presenter.BuildPaginationButtons(msg.ID, 0, totalPages)
	_, _ = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Components: &updatedComponents,
	})

	slog.Info("command completed", "command", "recommend",
		"track_name", output.SeedTrack.Name,
		"mode", output.Mode,
		"result_count", len(output.Items),
		"message_id", msg.ID)
}
