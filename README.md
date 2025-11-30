# jamberry 🍇

Spotify の楽曲・アーティスト・アルバム情報を Discord 上で簡単に検索・共有できる Bot です。

## 機能

### スラッシュコマンド

| コマンド                         | 説明                                      |
| -------------------------------- | ----------------------------------------- |
| `/track <spotify_url_or_id>`     | 楽曲情報を表示（KKBOX リンク付き）        |
| `/artist <spotify_url_or_id>`    | アーティスト情報を表示                    |
| `/album <spotify_url_or_id>`     | アルバム情報を表示                        |
| `/recommend <spotify_url_or_id>` | 楽曲に基づくレコメンドを表示（5 件）      |
| `/search <query>`                | 楽曲を検索（10 件、ページネーション対応） |

### 対応する Spotify ID 形式

- **Spotify URL**: `https://open.spotify.com/track/xxxxx`
- **Spotify URI**: `spotify:track:xxxxx`
- **Spotify ID**: `xxxxx`（22 文字の英数字）

## アーキテクチャ

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Discord   │────▶│  jamberry   │────▶│ tracktaste  │
│   (User)    │◀────│   (Bot)     │◀────│   (API)     │
└─────────────┘     └──────┬──────┘     └─────────────┘
                           │
                    ┌──────┴──────┐
                    │             │
               ┌────▼────┐  ┌─────▼─────┐
               │ L1 Cache │  │ L2 Cache  │
               │ (Memory) │  │  (Redis)  │
               └──────────┘  └───────────┘
```

### キャッシュ戦略

- **L1 キャッシュ**: インメモリ（sync.Map）、TTL 10 分
- **L2 キャッシュ**: Redis、TTL 30 日

### レート制限

- ユーザーごとに 10 秒間で最大 5 リクエスト（スライディングウィンドウ方式）

## セットアップ

### 前提条件

- Go 1.23+
- Docker & Docker Compose（推奨）
- Discord Bot Token
- Spotify Developer アカウント（Client ID / Secret）
- KKBOX Developer アカウント（Client ID / Secret）

### 環境変数

`.env` ファイルを作成してください：

```bash
cp .env.example .env
# .env ファイルを編集して必要な値を設定
```

| 変数名                  | 説明                               | 必須             |
| ----------------------- | ---------------------------------- | ---------------- |
| `DISCORD_BOT_TOKEN`     | Discord Bot のトークン             | ✅               |
| `SPOTIFY_CLIENT_ID`     | Spotify API の Client ID           | ✅               |
| `SPOTIFY_CLIENT_SECRET` | Spotify API の Client Secret       | ✅               |
| `KKBOX_ID`              | KKBOX API の Client ID             | ✅               |
| `KKBOX_SECRET`          | KKBOX API の Client Secret         | ✅               |
| `LOG_LEVEL`             | ログレベル (debug/info/warn/error) | デフォルト: info |

### Discord Bot の設定

1. [Discord Developer Portal](https://discord.com/developers/applications) でアプリケーションを作成
2. Bot を追加し、Token を取得
3. OAuth2 > URL Generator で以下のスコープを選択：
   - `bot`
   - `applications.commands`
4. Bot Permissions で以下を選択：
   - Send Messages
   - Embed Links
   - Use Slash Commands
5. 生成された URL でサーバーに Bot を招待

#### 招待 URL 例

```
https://discord.com/api/oauth2/authorize?client_id=YOUR_CLIENT_ID&permissions=2147485696&scope=bot%20applications.commands
```

### 起動方法

#### Docker Compose（推奨）

jamberry、tracktaste、Redis がすべて起動します。

```bash
# 起動
docker compose up -d

# ログ確認
docker compose logs -f jamberry

# 停止
docker compose down
```

#### ローカル実行

```bash
# Redis を起動
docker run -d -p 6379:6379 redis:7-alpine

# tracktaste を起動（別ターミナル）
# tracktaste の環境変数を設定して起動

# アプリケーションを起動
go run ./cmd/jamberry
```

## 開発

### プロジェクト構成

```
jamberry/
├── cmd/
│   └── jamberry/
│       └── main.go          # エントリーポイント
├── internal/
│   ├── cache/
│   │   └── cache.go         # 2層キャッシュ管理
│   ├── config/
│   │   └── config.go        # 設定読み込み
│   ├── embed/
│   │   └── builder.go       # Discord Embed 構築
│   ├── handler/
│   │   └── handler.go       # コマンド・ボタンハンドラー
│   ├── logger/
│   │   └── logger.go        # 構造化ロギング
│   ├── ratelimit/
│   │   └── limiter.go       # レート制限
│   ├── spotify/
│   │   └── validator.go     # Spotify ID 検証・正規化
│   └── tracktaste/
│       └── client.go        # tracktaste API クライアント
├── docs/
│   └── spec/
│       ├── SPEC.md          # 技術仕様書
│       └── USECASE.md       # ユースケース仕様書
├── compose.yml
├── Dockerfile
├── go.mod
└── README.md
```

### ビルド

```bash
go build -o jamberry ./cmd/jamberry
```

### テスト

```bash
go test ./...
```

## 依存サービス

### tracktaste

Spotify および KKBOX の情報を取得するバックエンド API。

- リポジトリ: `ghcr.io/t1nyb0x/tracktaste`
- 内部ポート: 8080

### Redis

L2 キャッシュとして使用。

- イメージ: `redis:7-alpine`
- 内部ポート: 6379

## ライセンス

MIT License - 詳細は [LICENSE](LICENSE) を参照してください。
