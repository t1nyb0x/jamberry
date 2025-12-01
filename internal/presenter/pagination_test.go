package presenter

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/t1nyb0x/jamberry/internal/domain"
)

func TestBuildRecommendEmbed(t *testing.T) {
	// テスト用のSimilarTrackを作成
	items := createTestSimilarTracks(30)

	tests := []struct {
		name              string
		originalTrackName string
		items             []domain.SimilarTrack
		page              int
		pageSize          int
		total             int
		wantTitle         string
		wantItemsInDesc   int
		wantStartNum      int
		wantEndNum        int
	}{
		{
			name:              "first page",
			originalTrackName: "Test Song",
			items:             items,
			page:              0,
			pageSize:          5,
			total:             30,
			wantTitle:         "🎶 おすすめトラック",
			wantItemsInDesc:   5,
			wantStartNum:      1,
			wantEndNum:        5,
		},
		{
			name:              "second page",
			originalTrackName: "Test Song",
			items:             items,
			page:              1,
			pageSize:          5,
			total:             30,
			wantTitle:         "🎶 おすすめトラック",
			wantItemsInDesc:   5,
			wantStartNum:      6,
			wantEndNum:        10,
		},
		{
			name:              "last page with less items",
			originalTrackName: "Test Song",
			items:             items[:7], // 7アイテムのみ
			page:              1,
			pageSize:          5,
			total:             7,
			wantTitle:         "🎶 おすすめトラック",
			wantItemsInDesc:   2, // 6, 7のみ
			wantStartNum:      6,
			wantEndNum:        7,
		},
		{
			name:              "empty items",
			originalTrackName: "Test Song",
			items:             []domain.SimilarTrack{},
			page:              0,
			pageSize:          5,
			total:             0,
			wantTitle:         "🎶 おすすめトラック",
			wantItemsInDesc:   0,
			wantStartNum:      1,
			wantEndNum:        0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			embed := BuildRecommendEmbed(tt.originalTrackName, tt.items, tt.page, tt.pageSize, tt.total, domain.RecommendModeBalanced)

			if embed.Title != tt.wantTitle {
				t.Errorf("Title = %s, want %s", embed.Title, tt.wantTitle)
			}

			if embed.Color != SpotifyGreen {
				t.Errorf("Color = %d, want %d", embed.Color, SpotifyGreen)
			}

			// Descriptionにオリジナルトラック名が含まれていることを確認
			if !strings.Contains(embed.Description, tt.originalTrackName) {
				t.Errorf("Description should contain original track name %s", tt.originalTrackName)
			}

			// 表示件数の確認（Descriptionの中のトラックリストを数える）
			if tt.wantItemsInDesc > 0 {
				// 番号付きリストをカウント
				for i := tt.wantStartNum; i <= tt.wantEndNum; i++ {
					expectedNum := fmt.Sprintf("**%d.", i)
					if !strings.Contains(embed.Description, expectedNum) {
						t.Errorf("Description should contain item number %d", i)
					}
				}
			}
		})
	}
}

func TestBuildSearchEmbed(t *testing.T) {
	// テスト用のTrackを作成
	items := createTestTracks(30)

	tests := []struct {
		name            string
		query           string
		items           []domain.Track
		page            int
		pageSize        int
		total           int
		wantTitle       string
		wantItemsInDesc int
		wantStartNum    int
		wantEndNum      int
	}{
		{
			name:            "first page",
			query:           "米津玄師",
			items:           items,
			page:            0,
			pageSize:        5,
			total:           30,
			wantTitle:       "🔍 検索結果",
			wantItemsInDesc: 5,
			wantStartNum:    1,
			wantEndNum:      5,
		},
		{
			name:            "second page",
			query:           "米津玄師",
			items:           items,
			page:            1,
			pageSize:        5,
			total:           30,
			wantTitle:       "🔍 検索結果",
			wantItemsInDesc: 5,
			wantStartNum:    6,
			wantEndNum:      10,
		},
		{
			name:            "last page with less items",
			query:           "YOASOBI",
			items:           items[:8],
			page:            1,
			pageSize:        5,
			total:           8,
			wantTitle:       "🔍 検索結果",
			wantItemsInDesc: 3, // 6, 7, 8のみ
			wantStartNum:    6,
			wantEndNum:      8,
		},
		{
			name:            "empty results",
			query:           "存在しないアーティスト",
			items:           []domain.Track{},
			page:            0,
			pageSize:        5,
			total:           0,
			wantTitle:       "🔍 検索結果",
			wantItemsInDesc: 0,
			wantStartNum:    1,
			wantEndNum:      0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			embed := BuildSearchEmbed(tt.query, tt.items, tt.page, tt.pageSize, tt.total)

			if embed.Title != tt.wantTitle {
				t.Errorf("Title = %s, want %s", embed.Title, tt.wantTitle)
			}

			if embed.Color != SpotifyGreen {
				t.Errorf("Color = %d, want %d", embed.Color, SpotifyGreen)
			}

			// Descriptionにクエリが含まれていることを確認
			if !strings.Contains(embed.Description, tt.query) {
				t.Errorf("Description should contain query %s", tt.query)
			}

			// 表示件数の確認
			if tt.wantItemsInDesc > 0 {
				for i := tt.wantStartNum; i <= tt.wantEndNum; i++ {
					expectedNum := fmt.Sprintf("**%d.", i)
					if !strings.Contains(embed.Description, expectedNum) {
						t.Errorf("Description should contain item number %d", i)
					}
				}
			}
		})
	}
}

func TestBuildPaginationButtons(t *testing.T) {
	tests := []struct {
		name           string
		messageID      string
		page           int
		totalPages     int
		wantPrevDisabled bool
		wantNextDisabled bool
	}{
		{
			name:             "first page",
			messageID:        "msg123",
			page:             0,
			totalPages:       6,
			wantPrevDisabled: true,
			wantNextDisabled: false,
		},
		{
			name:             "middle page",
			messageID:        "msg123",
			page:             2,
			totalPages:       6,
			wantPrevDisabled: false,
			wantNextDisabled: false,
		},
		{
			name:             "last page",
			messageID:        "msg123",
			page:             5,
			totalPages:       6,
			wantPrevDisabled: false,
			wantNextDisabled: true,
		},
		{
			name:             "single page",
			messageID:        "msg123",
			page:             0,
			totalPages:       1,
			wantPrevDisabled: true,
			wantNextDisabled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			components := BuildPaginationButtons(tt.messageID, tt.page, tt.totalPages)

			if len(components) != 1 {
				t.Fatalf("Expected 1 component row, got %d", len(components))
			}

			row, ok := components[0].(discordgo.ActionsRow)
			if !ok {
				t.Fatal("Expected ActionsRow component")
			}

			if len(row.Components) != 3 {
				t.Fatalf("Expected 3 buttons, got %d", len(row.Components))
			}

			// 前へボタン
			prevBtn, ok := row.Components[0].(discordgo.Button)
			if !ok {
				t.Fatal("Expected Button component for prev")
			}
			if prevBtn.Label != "◀ 前へ" {
				t.Errorf("Prev button label = %s, want ◀ 前へ", prevBtn.Label)
			}
			if prevBtn.Disabled != tt.wantPrevDisabled {
				t.Errorf("Prev button disabled = %v, want %v", prevBtn.Disabled, tt.wantPrevDisabled)
			}
			expectedPrevID := fmt.Sprintf("page_prev:%s:%d", tt.messageID, tt.page)
			if prevBtn.CustomID != expectedPrevID {
				t.Errorf("Prev button CustomID = %s, want %s", prevBtn.CustomID, expectedPrevID)
			}

			// 次へボタン
			nextBtn, ok := row.Components[1].(discordgo.Button)
			if !ok {
				t.Fatal("Expected Button component for next")
			}
			if nextBtn.Label != "次へ ▶" {
				t.Errorf("Next button label = %s, want 次へ ▶", nextBtn.Label)
			}
			if nextBtn.Disabled != tt.wantNextDisabled {
				t.Errorf("Next button disabled = %v, want %v", nextBtn.Disabled, tt.wantNextDisabled)
			}
			expectedNextID := fmt.Sprintf("page_next:%s:%d", tt.messageID, tt.page)
			if nextBtn.CustomID != expectedNextID {
				t.Errorf("Next button CustomID = %s, want %s", nextBtn.CustomID, expectedNextID)
			}

			// 自分も見るボタン
			viewBtn, ok := row.Components[2].(discordgo.Button)
			if !ok {
				t.Fatal("Expected Button component for view_own")
			}
			if viewBtn.Label != "👁 自分も見る" {
				t.Errorf("View button label = %s, want 👁 自分も見る", viewBtn.Label)
			}
			expectedViewID := fmt.Sprintf("view_own:%s", tt.messageID)
			if viewBtn.CustomID != expectedViewID {
				t.Errorf("View button CustomID = %s, want %s", viewBtn.CustomID, expectedViewID)
			}
			if viewBtn.Style != discordgo.PrimaryButton {
				t.Errorf("View button style = %v, want PrimaryButton", viewBtn.Style)
			}
		})
	}
}

// createTestSimilarTracks はテスト用のSimilarTrackを作成します
func createTestSimilarTracks(count int) []domain.SimilarTrack {
	tracks := make([]domain.SimilarTrack, count)
	for i := 0; i < count; i++ {
		tracks[i] = domain.SimilarTrack{
			ID:   fmt.Sprintf("track%d", i+1),
			Name: fmt.Sprintf("Track %d", i+1),
			URL:  fmt.Sprintf("https://open.spotify.com/track/track%d", i+1),
			Album: domain.Album{
				Name: fmt.Sprintf("Album %d", i+1),
				Artists: []domain.Artist{
					{Name: fmt.Sprintf("Artist %d", i+1)},
				},
			},
		}
	}
	return tracks
}

// createTestTracks はテスト用のTrackを作成します
func createTestTracks(count int) []domain.Track {
	tracks := make([]domain.Track, count)
	for i := 0; i < count; i++ {
		tracks[i] = domain.Track{
			ID:   fmt.Sprintf("track%d", i+1),
			Name: fmt.Sprintf("Track %d", i+1),
			URL:  fmt.Sprintf("https://open.spotify.com/track/track%d", i+1),
			Album: domain.Album{
				Name: fmt.Sprintf("Album %d", i+1),
			},
			Artists: []domain.Artist{
				{Name: fmt.Sprintf("Artist %d", i+1)},
			},
		}
	}
	return tracks
}
