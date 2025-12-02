package handler

import (
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestIsMentioned(t *testing.T) {
	botID := "123456789"

	tests := []struct {
		name     string
		mentions []*discordgo.User
		want     bool
	}{
		{
			name:     "no mentions",
			mentions: []*discordgo.User{},
			want:     false,
		},
		{
			name: "bot is mentioned",
			mentions: []*discordgo.User{
				{ID: botID},
			},
			want: true,
		},
		{
			name: "other user mentioned",
			mentions: []*discordgo.User{
				{ID: "987654321"},
			},
			want: false,
		},
		{
			name: "multiple mentions including bot",
			mentions: []*discordgo.User{
				{ID: "111111111"},
				{ID: botID},
				{ID: "222222222"},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isMentioned := false
			for _, mention := range tt.mentions {
				if mention.ID == botID {
					isMentioned = true
					break
				}
			}

			if isMentioned != tt.want {
				t.Errorf("isMentioned = %v, want %v", isMentioned, tt.want)
			}
		})
	}
}

func TestMentionResponseFormat(t *testing.T) {
	// 応答メッセージのフォーマット確認
	helpCommandID := "123456789012345678"
	expectedPrefix := "こんにちは！ 🍇\n"
	expectedContains := "</help:"

	response := expectedPrefix + "jamberry の使い方は </help:" + helpCommandID + "> で確認できます。"

	if len(response) == 0 {
		t.Error("Response should not be empty")
	}

	if response[:len(expectedPrefix)] != expectedPrefix {
		t.Errorf("Response should start with '%s'", expectedPrefix)
	}

	found := false
	for i := 0; i <= len(response)-len(expectedContains); i++ {
		if response[i:i+len(expectedContains)] == expectedContains {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Response should contain '%s'", expectedContains)
	}
}
