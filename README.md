# jamberry

Spotify 関連の情報を Discord 上で取得できる Discord Bot です。

## 機能

- `/track <url>` - トラック情報の取得
- `/artist <url>` - アーティスト情報の取得
- `/album <url>` - アルバム情報の取得
- `/recommend <url>` - トラックに基づくレコメンド取得
- `/search <query>` - トラック検索

## セットアップ

### 前提条件

- Go 1.23+
- Docker & Docker Compose (推奨)
- Discord Bot Token
- Spotify Developer アカウント（Client ID / Secret）
- KKBOX Developer アカウント（Client ID / Secret）

### 環境変数

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
| `LOG_LEVEL`             | ログレベル (DEBUG/INFO/WARN/ERROR) | デフォルト: INFO |

### 実行方法

#### Docker Compose (推奨)

jamberry、tracktaste、Redis がすべて起動します。

```bash
docker compose up -d
```

#### ローカル実行

```bash
# Redis を起動
docker run -d -p 6379:6379 redis:7-alpine

# アプリケーションを起動
go run ./cmd/jamberry
```

## 開発

### ビルド

```bash
go build -o jamberry ./cmd/jamberry
```

### テスト

```bash
go test ./...
```

## Discord Bot の設定

### 必要な権限

- Send Messages
- Embed Links
- Use Slash Commands

### OAuth2 スコープ

- `bot`
- `applications.commands`

### 招待 URL

```
https://discord.com/api/oauth2/authorize?client_id=YOUR_CLIENT_ID&permissions=2147485696&scope=bot%20applications.commands
```

## ライセンス

LICENSE ファイルを参照してください。
