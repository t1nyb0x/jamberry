# アーキテクチャ仕様書

## 概要

jamberry は**レイヤードアーキテクチャ**（Layered Architecture）を採用した Discord Bot です。
各レイヤーは明確な責務を持ち、依存性の方向は内側（domain）に向かいます。

## レイヤー構成

```
┌─────────────────────────────────────────────────────────────┐
│                      cmd/server/main.go                     │
│                      (エントリポイント)                      │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                          bot 層                              │
│                  Discord Bot セッション管理                   │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                       handler 層                             │
│              Discord インタラクション処理                     │
│         (コマンド受信 → usecase 呼び出し → 応答)              │
└─────────────────────────────────────────────────────────────┘
                              │
              ┌───────────────┼───────────────┐
              ▼               ▼               ▼
┌───────────────────┐ ┌─────────────┐ ┌─────────────────────┐
│    usecase 層      │ │ presenter層 │ │    ratelimit 層     │
│  ビジネスロジック   │ │ Embed 構築  │ │    レート制限        │
└───────────────────┘ └─────────────┘ └─────────────────────┘
              │
              ▼
┌─────────────────────────────────────────────────────────────┐
│                        domain 層                             │
│            エンティティ・リポジトリインターフェース             │
└─────────────────────────────────────────────────────────────┘
              ▲
              │ (インターフェース実装)
              │
┌─────────────────────────────────────────────────────────────┐
│                    infrastructure 層                         │
│          外部サービス連携（tracktaste, Redis）                │
└─────────────────────────────────────────────────────────────┘
```

## 各レイヤーの責務

### 1. domain 層 (`internal/domain/`)

**責務**: ビジネスエンティティとリポジトリインターフェースの定義

```
internal/domain/
├── track.go       # Track, SimilarTrack エンティティ
├── artist.go      # Artist, ArtistDetail エンティティ
├── album.go       # Album, AlbumDetail, AlbumTrack, Image エンティティ
├── cache.go       # PaginationData, CacheRepository インターフェース
└── repository.go  # TrackRepository, ArtistRepository, AlbumRepository, MusicRepository
```

**特徴**:

- 外部依存なし（標準ライブラリのみ）
- 他のレイヤーから参照される
- ビジネスルールをカプセル化

**インターフェース定義**:

```go
// MusicRepository は音楽情報を取得する統合リポジトリインターフェース
type MusicRepository interface {
    TrackRepository
    ArtistRepository
    AlbumRepository
}

// CacheRepository はキャッシュを操作するリポジトリインターフェース
type CacheRepository interface {
    Set(ctx context.Context, key string, data *PaginationData) error
    Get(ctx context.Context, key string) (*PaginationData, error)
    Delete(ctx context.Context, key string)
}
```

### 2. usecase 層 (`internal/usecase/`)

**責務**: ビジネスロジックの実装

```
internal/usecase/
├── track.go      # TrackUseCase - トラック情報取得
├── artist.go     # ArtistUseCase - アーティスト情報取得
├── album.go      # AlbumUseCase - アルバム情報取得
├── recommend.go  # RecommendUseCase - レコメンド取得
├── search.go     # SearchUseCase - トラック検索
└── errors.go     # ValidationError, NotFoundError
```

**特徴**:

- domain 層のインターフェースに依存
- 入力バリデーション
- ビジネスロジックの実行
- 適切なエラーハンドリング

**パターン**:

```go
type TrackUseCase struct {
    repo domain.TrackRepository  // インターフェースに依存
}

func (u *TrackUseCase) GetTrack(ctx context.Context, input TrackInput) (*TrackOutput, error) {
    // 1. バリデーション
    // 2. リポジトリ呼び出し
    // 3. 結果返却
}
```

### 3. handler 層 (`internal/handler/`)

**責務**: Discord インタラクションの処理

```
internal/handler/
├── handler.go     # ルーター（コマンド振り分け）
├── track.go       # /jam track コマンドハンドラー
├── artist.go      # /jam artist コマンドハンドラー
├── album.go       # /jam album コマンドハンドラー
├── recommend.go   # /jam recommend コマンドハンドラー
├── search.go      # /jam search コマンドハンドラー
├── component.go   # ボタンインタラクションハンドラー
└── responder.go   # Discord レスポンスヘルパー
```

**特徴**:

- Discord 固有のロジックを担当
- usecase 層を呼び出し
- presenter 層で Embed を構築
- レスポンス送信

**フロー**:

```
Discord → handler.HandleInteraction()
                    │
                    ├─ handleCommand() ─→ handleTrack/Artist/Album/...
                    │                            │
                    │                            ▼
                    │                      usecase.GetXxx()
                    │                            │
                    │                            ▼
                    │                      presenter.BuildXxxEmbed()
                    │                            │
                    │                            ▼
                    │                      responder.EditResponseEmbed()
                    │
                    └─ handleComponent() ─→ handlePaging/handleViewOwn
```

### 4. presenter 層 (`internal/presenter/`)

**責務**: Discord Embed の構築

```
internal/presenter/
├── embed.go       # Track/Artist/Album Embed 構築
├── pagination.go  # Recommend/Search Embed + ページネーションボタン
└── formatter.go   # FormatDuration, FormatNumber, GetLargestImage, JoinArtistNames
```

**特徴**:

- domain エンティティを Discord Embed に変換
- 表示ロジックをカプセル化
- 再利用可能なフォーマッター

### 5. infrastructure 層 (`internal/infrastructure/`)

**責務**: 外部サービスとの連携

```
internal/infrastructure/
├── tracktaste/           # tracktaste API クライアント
│   ├── client.go         # HTTP クライアント、domain.MusicRepository 実装
│   ├── track.go          # Track レスポンス → domain.Track 変換
│   ├── artist.go         # Artist レスポンス → domain.ArtistDetail 変換
│   └── album.go          # Album レスポンス → domain.AlbumDetail 変換
└── cache/
    └── cache.go          # L1/L2 キャッシュ、domain.CacheRepository 実装
```

**特徴**:

- domain 層のインターフェースを実装
- 外部 API レスポンスを domain エンティティに変換
- エラーハンドリングとログ出力

**インターフェース実装の確認**:

```go
// コンパイル時にインターフェース実装を検証
var _ domain.MusicRepository = (*Client)(nil)
var _ domain.CacheRepository = (*Manager)(nil)
```

### 6. bot 層 (`internal/bot/`)

**責務**: Discord Bot セッション管理

```
internal/bot/
├── bot.go        # Bot 構造体、Start/Stop、ハンドラー登録
└── commands.go   # スラッシュコマンド定義
```

### 7. その他のパッケージ

| パッケージ  | 責務                                 |
| ----------- | ------------------------------------ |
| `config`    | 環境変数からの設定読み込み           |
| `logger`    | 構造化ロギング（slog）のセットアップ |
| `ratelimit` | ユーザーごとのレート制限             |
| `spotify`   | Spotify URL/URI/ID のバリデーション  |

## 依存関係

```
main.go
    │
    ├── bot
    │     └── (discordgo)
    │
    ├── handler
    │     ├── usecase
    │     │     ├── domain
    │     │     └── spotify
    │     ├── presenter
    │     │     └── domain
    │     └── ratelimit
    │
    ├── infrastructure/tracktaste
    │     └── domain
    │
    ├── infrastructure/cache
    │     └── domain
    │
    ├── config
    └── logger
```

**依存性の逆転 (Dependency Inversion)**:

- usecase 層は domain 層のインターフェース（`MusicRepository`）に依存
- infrastructure 層がそのインターフェースを実装
- main.go で具体的な実装を注入（Dependency Injection）

```go
// main.go での依存性注入
ttClient := tracktaste.NewClient(cfg.TrackTasteAPIURL)  // 具体的な実装
trackUC := usecase.NewTrackUseCase(ttClient)             // インターフェースとして渡す
```

## データフロー

### コマンド実行フロー（例: /jam track）

```
1. Discord → InteractionCreate イベント
       │
2.     ▼
   handler.HandleInteraction()
       │
3.     ▼
   handleCommand() → handleTrack()
       │
4.     ▼
   trackUseCase.GetTrack(input)
       │
5.     ├─→ spotify.ValidateInput() [バリデーション]
       │
6.     └─→ repo.FetchTrack() [tracktaste API 呼び出し]
              │
7.            ▼
         tracktaste.Client.FetchTrack()
              │
8.            ▼
         HTTP GET /v1/track/fetch?url=...
              │
9.            ▼
         trackResponse.toDomain() [domain.Track に変換]
              │
10.           ▼
         presenter.BuildTrackEmbed(track)
              │
11.           ▼
         responder.EditResponseEmbed(embed)
              │
12.           ▼
         Discord へ応答
```

### ページネーションフロー

```
1. /jam recommend or /jam search コマンド
       │
2.     ▼
   結果を cache に保存（key: messageID）
       │
3.     ▼
   ボタン付き Embed を送信
       │
4. ユーザーがボタンをクリック
       │
5.     ▼
   handler.handleComponent()
       │
6.     ├─→ cache.Get(messageID) [キャッシュからデータ取得]
       │
7.     └─→ 権限チェック（OwnerID 比較）
              │
8.            ▼
         presenter.BuildXxxEmbed(items, newPage)
              │
9.            ▼
         responder.UpdateMessage(embed, components)
```

## 設計原則

### 1. 単一責任の原則 (SRP)

各パッケージ・ファイルは単一の責務を持つ:

- `handler/track.go` → /jam track コマンドの処理のみ
- `presenter/embed.go` → Embed 構築のみ
- `infrastructure/tracktaste/client.go` → API 通信のみ

### 2. 開放閉鎖の原則 (OCP)

新しいコマンドを追加する場合:

1. `usecase/` に新しいユースケースを追加
2. `handler/` に新しいハンドラーを追加
3. `bot/commands.go` にコマンド定義を追加
4. 既存コードを変更せずに拡張可能

### 3. 依存性逆転の原則 (DIP)

- 上位レイヤー（usecase）は下位レイヤー（infrastructure）に直接依存しない
- インターフェース（domain）を介して依存

### 4. テスト容易性

インターフェースによりモック可能:

```go
// テスト用モック
type mockTrackRepo struct {
    track *domain.Track
    err   error
}

func (m *mockTrackRepo) FetchTrack(ctx context.Context, url string) (*domain.Track, error) {
    return m.track, m.err
}

// テストコード
func TestGetTrack(t *testing.T) {
    mockRepo := &mockTrackRepo{track: &domain.Track{Name: "Test"}}
    uc := usecase.NewTrackUseCase(mockRepo)
    // ...
}
```

## ファイル一覧

```
internal/
├── bot/
│   ├── bot.go
│   └── commands.go
├── config/
│   └── config.go
├── domain/
│   ├── album.go
│   ├── artist.go
│   ├── cache.go
│   ├── repository.go
│   └── track.go
├── handler/
│   ├── album.go
│   ├── artist.go
│   ├── component.go
│   ├── handler.go
│   ├── recommend.go
│   ├── responder.go
│   ├── search.go
│   └── track.go
├── infrastructure/
│   ├── cache/
│   │   └── cache.go
│   └── tracktaste/
│       ├── album.go
│       ├── artist.go
│       ├── client.go
│       └── track.go
├── logger/
│   └── logger.go
├── presenter/
│   ├── embed.go
│   ├── formatter.go
│   └── pagination.go
├── ratelimit/
│   └── limiter.go
├── spotify/
│   └── validator.go
└── usecase/
    ├── album.go
    ├── artist.go
    ├── errors.go
    ├── recommend.go
    ├── search.go
    └── track.go
```

---

## 関連ドキュメント

| ドキュメント                   | 説明                                               |
| ------------------------------ | -------------------------------------------------- |
| [SPEC.md](./SPEC.md)           | 技術仕様（機能一覧、API インターフェース、エラー） |
| [USECASE.md](./USECASE.md)     | ユースケース仕様（ユーザーストーリー、フロー）     |
| [TEST_SPEC.md](./TEST_SPEC.md) | テスト仕様（テストケース、カバレッジ目標）         |
