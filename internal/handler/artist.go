package handler

import (
	"context"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/t1nyb0x/jamberry/internal/presenter"
	"github.com/t1nyb0x/jamberry/internal/usecase"
)

// handleArtist はアーティスト情報取得コマンドを処理します
func (h *Handler) handleArtist(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) {
	if len(options) == 0 {
		slog.Info("validation failed: empty input", "command", "jam artist")
		h.responder.RespondEphemeral(s, i, "❌ URL を入力してください。")
		return
	}

	input := options[0].StringValue()

	// DeferReply
	if err := h.responder.DeferReply(s, i); err != nil {
		slog.Error("failed to defer reply", "command", "jam artist", "error", err)
		return
	}

	ctx := context.Background()
	output, err := h.artistUseCase.GetArtist(ctx, usecase.ArtistInput{Input: input})
	if err != nil {
		h.responder.EditResponse(s, i, err.Error())
		return
	}

	// Embed構築・返信
	emb := presenter.BuildArtistEmbed(output.Artist)
	h.responder.EditResponseEmbed(s, i, emb)
	slog.Info("command completed", "command", "jam artist", "artist_name", output.Artist.Name, "artist_id", output.Artist.ID)
}
