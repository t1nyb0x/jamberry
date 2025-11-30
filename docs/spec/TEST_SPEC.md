# jamberry テスト仕様書

## 概要

本ドキュメントは jamberry のユニットテスト仕様を定義します。テストは Go の標準テストフレームワークを使用して実装されています。

## テスト対象パッケージ

| パッケージ           | カバレッジ | 説明                                       |
| -------------------- | ---------- | ------------------------------------------ |
| `internal/spotify`   | 95.3%      | Spotify 入力バリデーション                 |
| `internal/presenter` | 100%       | Discord Embed 構築・フォーマット処理       |
| `internal/ratelimit` | 81.6%      | アプリケーションレベルのレートリミット処理 |

---

## 1. Spotify Validator テスト (`internal/spotify/validator_test.go`)

### 1.1 空入力テスト (`TestValidateInput_EmptyInput`)

| テストケース    | 入力値  | 期待結果                                    |
| --------------- | ------- | ------------------------------------------- |
| empty string    | `""`    | `Valid=false`, エラー: 認識できませんでした |
| whitespace only | `"   "` | `Valid=false`, エラー: 認識できませんでした |
| tab only        | `"\t"`  | `Valid=false`, エラー: 認識できませんでした |
| newline only    | `"\n"`  | `Valid=false`, エラー: 認識できませんでした |

### 1.2 有効な URL テスト (`TestValidateInput_ValidURL`)

| テストケース            | 入力値                                                            | 期待結果                           |
| ----------------------- | ----------------------------------------------------------------- | ---------------------------------- |
| track URL               | `https://open.spotify.com/track/4iV5W9uYEdYUVa79Axb7Rh`           | `Valid=true`, URL 正規化成功       |
| artist URL              | `https://open.spotify.com/artist/0OdUWJ0sBjDrqHygGUXeCF`          | `Valid=true`, URL 正規化成功       |
| album URL               | `https://open.spotify.com/album/4aawyAB9vmqN3uQ7FjRGTy`           | `Valid=true`, URL 正規化成功       |
| URL with query params   | `https://open.spotify.com/track/4iV5W9uYEdYUVa79Axb7Rh?si=abc123` | `Valid=true`, クエリパラメータ除去 |
| URL with trailing slash | `https://open.spotify.com/track/4iV5W9uYEdYUVa79Axb7Rh/`          | `Valid=true`, 末尾スラッシュ処理   |
| URL with whitespace     | ` https://open.spotify.com/track/4iV5W9uYEdYUVa79Axb7Rh `         | `Valid=true`, 前後空白除去         |

### 1.3 無効な URL テスト (`TestValidateInput_InvalidURL`)

| テストケース       | 入力値                                                   | 期待結果                                    |
| ------------------ | -------------------------------------------------------- | ------------------------------------------- |
| non-spotify domain | `https://example.com/track/4iV5W9uYEdYUVa79Axb7Rh`       | `Valid=false`, エラー: 認識できませんでした |
| youtube URL        | `https://www.youtube.com/watch?v=dQw4w9WgXcQ`            | `Valid=false`, エラー: 認識できませんでした |
| spotify embed URL  | `https://embed.spotify.com/track/4iV5W9uYEdYUVa79Axb7Rh` | `Valid=false`, エラー: 認識できませんでした |
| missing path       | `https://open.spotify.com/`                              | `Valid=false`, エラー: 認識できませんでした |
| invalid ID format  | `https://open.spotify.com/track/abc`                     | `Valid=false`, エラー: 認識できませんでした |

### 1.4 有効な URI テスト (`TestValidateInput_ValidURI`)

| テストケース        | 入力値                                   | 期待結果                   |
| ------------------- | ---------------------------------------- | -------------------------- |
| track URI           | `spotify:track:4iV5W9uYEdYUVa79Axb7Rh`   | `Valid=true`, URL に正規化 |
| artist URI          | `spotify:artist:0OdUWJ0sBjDrqHygGUXeCF`  | `Valid=true`, URL に正規化 |
| album URI           | `spotify:album:4aawyAB9vmqN3uQ7FjRGTy`   | `Valid=true`, URL に正規化 |
| URI with whitespace | ` spotify:track:4iV5W9uYEdYUVa79Axb7Rh ` | `Valid=true`, 前後空白除去 |

### 1.5 無効な URI テスト (`TestValidateInput_InvalidURI`)

| テストケース       | 入力値                                       | 期待結果                                    |
| ------------------ | -------------------------------------------- | ------------------------------------------- |
| not spotify prefix | `other:track:4iV5W9uYEdYUVa79Axb7Rh`         | `Valid=false`, エラー: 認識できませんでした |
| too few segments   | `spotify:track`                              | `Valid=false`, エラー: 認識できませんでした |
| too many segments  | `spotify:track:4iV5W9uYEdYUVa79Axb7Rh:extra` | `Valid=false`, エラー: 認識できませんでした |
| invalid ID in URI  | `spotify:track:abc`                          | `Valid=false`, エラー: 認識できませんでした |

### 1.6 有効な生 ID テスト (`TestValidateInput_ValidRawID`)

| テストケース              | 入力値                     | 期待結果                   |
| ------------------------- | -------------------------- | -------------------------- |
| valid 22-char ID (track)  | `4iV5W9uYEdYUVa79Axb7Rh`   | `Valid=true`, URL に正規化 |
| valid 22-char ID (artist) | `0OdUWJ0sBjDrqHygGUXeCF`   | `Valid=true`, URL に正規化 |
| valid 22-char ID (album)  | `4aawyAB9vmqN3uQ7FjRGTy`   | `Valid=true`, URL に正規化 |
| ID with whitespace        | ` 4iV5W9uYEdYUVa79Axb7Rh ` | `Valid=true`, 前後空白除去 |

### 1.7 無効な生 ID テスト (`TestValidateInput_InvalidRawID`)

| テストケース       | 入力値                      | 期待結果                                    |
| ------------------ | --------------------------- | ------------------------------------------- |
| too short ID       | `abc`                       | `Valid=false`, エラー: 認識できませんでした |
| too long ID        | `4iV5W9uYEdYUVa79Axb7Rh123` | `Valid=false`, エラー: 認識できませんでした |
| special characters | `4iV5W9uYEdYUVa79Axb7!@`    | `Valid=false`, エラー: 認識できませんでした |
| 21 characters      | `4iV5W9uYEdYUVa79Axb7R`     | `Valid=false`, エラー: 認識できませんでした |
| 23 characters      | `4iV5W9uYEdYUVa79Axb7Rha`   | `Valid=false`, エラー: 認識できませんでした |

### 1.8 エンティティ種別不一致テスト (`TestValidateInput_EntityTypeMismatch`)

| テストケース                   | 入力値 (URL/URI)                      | 期待種別 | 期待エラー                   |
| ------------------------------ | ------------------------------------- | -------- | ---------------------------- |
| track URL with artist expected | `https://open.spotify.com/track/xxx`  | artist   | ArtistURL を入力してください |
| track URL with album expected  | `https://open.spotify.com/track/xxx`  | album    | AlbumURL を入力してください  |
| artist URL with track expected | `https://open.spotify.com/artist/xxx` | track    | TrackURL を入力してください  |
| album URL with track expected  | `https://open.spotify.com/album/xxx`  | track    | TrackURL を入力してください  |
| track URI with artist expected | `spotify:track:xxx`                   | artist   | ArtistURL を入力してください |

### 1.9 エラーメッセージテスト (`TestGetEntityMismatchError`)

| エンティティ種別 | 期待メッセージ                                      |
| ---------------- | --------------------------------------------------- |
| track            | ❌ Spotify の TrackURL を入力してください           |
| artist           | ❌ Spotify の ArtistURL を入力してください          |
| album            | ❌ Spotify の AlbumURL を入力してください           |
| unknown          | ❌ Spotify の URL / ID として認識できませんでした。 |

### 1.10 ID 正規表現テスト (`TestSpotifyIDRegex`)

Spotify ID は 22 文字の英数字で構成される。

| 入力値                    | 期待結果 |
| ------------------------- | -------- |
| `4iV5W9uYEdYUVa79Axb7Rh`  | `true`   |
| `abc`                     | `false`  |
| `4iV5W9uYEdYUVa79Axb7Rh1` | `false`  |
| `4iV5W9uYEdYUVa79-xb7Rh`  | `false`  |
| `""`                      | `false`  |

---

## 2. Presenter Formatter テスト (`internal/presenter/formatter_test.go`)

### 2.1 再生時間フォーマットテスト (`TestFormatDuration`)

| テストケース          | 入力 (ms) | 期待出力 |
| --------------------- | --------- | -------- |
| zero                  | 0         | `0:00`   |
| one second            | 1000      | `0:01`   |
| 59 seconds            | 59000     | `0:59`   |
| one minute            | 60000     | `1:00`   |
| one minute one second | 61000     | `1:01`   |
| 3:45 (typical song)   | 225000    | `3:45`   |
| 10 minutes            | 600000    | `10:00`  |
| over hour             | 3661000   | `61:01`  |

### 2.2 数値フォーマットテスト (`TestFormatNumber`)

| テストケース      | 入力      | 期待出力      |
| ----------------- | --------- | ------------- |
| zero              | 0         | `0`           |
| single digit      | 5         | `5`           |
| three digits      | 123       | `123`         |
| four digits       | 1234      | `1,234`       |
| million           | 1000000   | `1,000,000`   |
| large number      | 123456789 | `123,456,789` |
| typical followers | 5823914   | `5,823,914`   |

### 2.3 最大画像取得テスト (`TestGetLargestImage`)

| テストケース                    | 画像データ                | 期待結果                  |
| ------------------------------- | ------------------------- | ------------------------- |
| empty images                    | `[]`                      | `""`                      |
| single image                    | 640x640                   | 該当 URL                  |
| multiple images - largest first | 640x640, 300x300, 64x64   | 640x640 の URL            |
| multiple images - largest last  | 64x64, 300x300, 640x640   | 640x640 の URL            |
| different aspect ratios         | 800x400, 640x640, 400x800 | 640x640 の URL (最大面積) |

### 2.4 アーティスト名結合テスト (`TestJoinArtistNames`)

| テストケース               | アーティスト名             | 期待出力             |
| -------------------------- | -------------------------- | -------------------- |
| empty artists              | `[]`                       | `""`                 |
| single artist              | `["Artist A"]`             | `Artist A`           |
| two artists                | `["Artist A", "Artist B"]` | `Artist A, Artist B` |
| three artists              | `["A", "B", "C"]`          | `A, B, C`            |
| artists with special chars | `["米津玄師", "YOASOBI"]`  | `米津玄師, YOASOBI`  |

---

## 3. Presenter Embed テスト (`internal/presenter/embed_test.go`)

### 3.1 トラック Embed テスト (`TestBuildTrackEmbed`)

| テストケース             | 入力条件         | 検証項目                            |
| ------------------------ | ---------------- | ----------------------------------- |
| basic track              | 全フィールド有り | タイトル、説明、URL、色、サムネイル |
| explicit track           | `Explicit=true`  | タイトルに 🔞 が付与される          |
| track without popularity | `Popularity=nil` | 人気度フィールドが省略される        |
| track without album art  | `Images=[]`      | サムネイルが設定されない            |

#### 検証内容

- タイトル形式: `🎵 {トラック名}` (Explicit の場合は `🎵 {トラック名} 🔞`)
- 説明: アーティスト名のカンマ区切り
- 色: SpotifyGreen (`0x1DB954`)
- フィールド: アルバム、再生時間、リリース日、(人気度)

### 3.2 アーティスト Embed テスト (`TestBuildArtistEmbed`)

| テストケース                   | 入力条件         | 検証項目                     |
| ------------------------------ | ---------------- | ---------------------------- |
| artist with all fields         | 全フィールド有り | タイトル、URL、サムネイル    |
| artist without genres          | `Genres=[]`      | ジャンル: `なし` と表示      |
| artist with more than 3 genres | ジャンル 5 件    | 最大 3 件まで表示            |
| artist without popularity      | `Popularity=nil` | 人気度フィールドが省略される |

#### 検証内容

- タイトル形式: `🎤 {アーティスト名}`
- フィールド: ジャンル（最大 3 件）、フォロワー、(人気度)

### 3.3 アルバム Embed テスト (`TestBuildAlbumEmbed`)

| テストケース                  | 入力条件         | 検証項目                           |
| ----------------------------- | ---------------- | ---------------------------------- |
| album with all fields         | 全フィールド有り | タイトル、説明、サムネイル、収録曲 |
| album without popularity      | `Popularity=nil` | 人気度フィールドが省略される       |
| album with no tracks          | `Tracks=[]`      | 収録曲フィールドが省略される       |
| album with exactly 5 tracks   | トラック 5 件    | 全 5 曲が表示される                |
| album with more than 5 tracks | トラック 7 件    | 先頭 5 曲のみ表示                  |

#### 検証内容

- タイトル形式: `💿 {アルバム名}`
- 説明: アーティスト名のカンマ区切り
- フィールド: リリース日、トラック数、(人気度)、収録曲（先頭 5 曲）

---

## 4. Presenter Pagination テスト (`internal/presenter/pagination_test.go`)

### 4.1 レコメンド Embed テスト (`TestBuildRecommendEmbed`)

| テストケース              | ページ | 総件数 | 期待表示件数 | 期待開始番号 |
| ------------------------- | ------ | ------ | ------------ | ------------ |
| first page                | 0      | 30     | 5            | 1            |
| second page               | 1      | 30     | 5            | 6            |
| last page with less items | 1      | 7      | 2            | 6            |
| empty items               | 0      | 0      | 0            | -            |

#### 検証内容

- タイトル: `🎶 おすすめトラック`
- 説明に元トラック名と件数情報が含まれる
- 各トラックにアーティスト名、アルバム名、Spotify リンクが表示される

### 4.2 検索 Embed テスト (`TestBuildSearchEmbed`)

| テストケース              | ページ | 総件数 | 期待表示件数 | 期待開始番号 |
| ------------------------- | ------ | ------ | ------------ | ------------ |
| first page                | 0      | 30     | 5            | 1            |
| second page               | 1      | 30     | 5            | 6            |
| last page with less items | 1      | 8      | 3            | 6            |
| empty results             | 0      | 0      | 0            | -            |

#### 検証内容

- タイトル: `🔍 検索結果`
- 説明に検索クエリと件数情報が含まれる

### 4.3 ページングボタンテスト (`TestBuildPaginationButtons`)

| テストケース | ページ | 総ページ数 | 前へ無効化 | 次へ無効化 |
| ------------ | ------ | ---------- | ---------- | ---------- |
| first page   | 0      | 6          | `true`     | `false`    |
| middle page  | 2      | 6          | `false`    | `false`    |
| last page    | 5      | 6          | `false`    | `true`     |
| single page  | 0      | 1          | `true`     | `true`     |

#### 検証内容

- ボタン 3 つ: `◀ 前へ`, `次へ ▶`, `👁 自分も見る`
- CustomID 形式:
  - 前へ: `page_prev:{messageID}:{page}`
  - 次へ: `page_next:{messageID}:{page}`
  - 自分も見る: `view_own:{messageID}`
- `👁 自分も見る` ボタンは常に PrimaryButton スタイル

---

## 5. Rate Limiter テスト (`internal/ratelimit/limiter_test.go`)

### 5.1 基本使用テスト (`TestLimiter_Allow_BasicUsage`)

| リクエスト回数 | 期待結果 |
| -------------- | -------- |
| 1-5 回目       | 許可     |
| 6 回目         | 拒否     |

### 5.2 異なるユーザーテスト (`TestLimiter_Allow_DifferentUsers`)

| ユーザー | リクエスト回数 | 期待結果 |
| -------- | -------------- | -------- |
| user1    | 1-5 回目       | 許可     |
| user1    | 6 回目         | 拒否     |
| user2    | 1-5 回目       | 許可     |
| user2    | 6 回目         | 拒否     |

→ ユーザーごとに独立してレートリミットが適用される

### 5.3 並行処理テスト (`TestLimiter_Allow_Concurrency`)

| 条件                       | 期待結果           |
| -------------------------- | ------------------ |
| 同一ユーザーで 10 並行実行 | 5 回許可、5 回拒否 |

→ 競合状態が発生しないことを確認

### 5.4 クリーンアップテスト (`TestLimiter_Cleanup`)

| テストケース           | 条件                      | 期待結果             |
| ---------------------- | ------------------------- | -------------------- |
| ウィンドウ内のエントリ | 10 秒以内のタイムスタンプ | エントリは保持される |
| 期限切れエントリ       | 20 秒前のタイムスタンプ   | エントリは削除される |

### 5.5 定数確認テスト (`TestConstants`)

| 定数名      | 期待値 |
| ----------- | ------ |
| Window      | 10 秒  |
| MaxRequests | 5      |

### 5.6 エッジケーステスト

| テストケース        | 条件           | 期待結果   |
| ------------------- | -------------- | ---------- |
| empty user ID       | `userID = ""`  | 正常動作   |
| exactly MaxRequests | 5 回リクエスト | 全て許可   |
| MaxRequests + 1     | 6 回リクエスト | 6 回目拒否 |

---

## テスト実行方法

### 全テスト実行

```bash
go test ./internal/... -v
```

### カバレッジ付き実行

```bash
go test ./internal/... -cover
```

### 特定パッケージのみ実行

```bash
# Spotify Validator
go test ./internal/spotify/... -v

# Presenter
go test ./internal/presenter/... -v

# Rate Limiter
go test ./internal/ratelimit/... -v
```

### カバレッジレポート生成

```bash
go test ./internal/... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

---

## カバレッジ目標

| パッケージ           | 目標 | 現状  | 備考    |
| -------------------- | ---- | ----- | ------- |
| `internal/spotify`   | 80%  | 95.3% | ✅ 達成 |
| `internal/presenter` | 80%  | 100%  | ✅ 達成 |
| `internal/ratelimit` | 80%  | 81.6% | ✅ 達成 |

### カバレッジが 100% に達しない箇所

#### `internal/ratelimit`

- `StartCleanup` 関数: ゴルーチンで定期実行されるため、テストでの検証が困難
- ただし、`Cleanup` 関数自体はテストでカバーされている

---

## 関連ドキュメント

- [jamberry 仕様書](./SPEC.md)
- [アーキテクチャ設計](./ARCHITECTURE.md)
- [ユースケース](./USECASE.md)
