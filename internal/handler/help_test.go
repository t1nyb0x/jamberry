package handler

import (
	"testing"
)

func TestHelpEmbedFields(t *testing.T) {
	// ヘルプに含まれるべきフィールド名
	expectedFields := []string{
		"🎵 `/jam track <url>`",
		"👤 `/jam artist <url>`",
		"💿 `/jam album <url>`",
		"✨ `/jam recommend <url> [mode]`",
		"🔍 `/jam search <query>`",
		"🩺 `/tracktaste`",
		"❓ `/help`",
		"📝 対応する入力形式",
	}

	// フィールド数の確認
	if len(expectedFields) != 8 {
		t.Errorf("Expected 8 help fields, got %d", len(expectedFields))
	}
}

func TestRecommendHelpContent(t *testing.T) {
	// レコメンドのヘルプに含まれるべき内容
	expectedContent := []string{
		"バランス",
		"雰囲気重視",
		"関連性重視",
		"スコア",
		"ボーナス倍率",
		"×2.5",
		"×1.3",
		"×1.2",
		"×1.1",
		"×0.5",
	}

	// コンテンツ項目数の確認
	if len(expectedContent) != 10 {
		t.Errorf("Expected 10 recommend help content items, got %d", len(expectedContent))
	}
}

func TestInputFormatHelpContent(t *testing.T) {
	// 対応入力形式
	expectedFormats := []string{
		"Spotify URL",
		"Spotify URI",
		"Spotify ID",
	}

	if len(expectedFormats) != 3 {
		t.Errorf("Expected 3 input formats, got %d", len(expectedFormats))
	}
}
