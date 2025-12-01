package version

// Version はアプリケーションのバージョン情報です
// ビルド時に -ldflags で上書き可能です
var (
	// Version はセマンティックバージョン
	Version = "0.2.0"

	// GitCommit はGitコミットハッシュ（ビルド時に設定）
	GitCommit = "unknown"

	// BuildDate はビルド日時（ビルド時に設定）
	BuildDate = "unknown"
)

// GetVersion はバージョン文字列を返します
func GetVersion() string {
	return Version
}

// GetFullVersion は完全なバージョン情報を返します
func GetFullVersion() string {
	return Version + " (" + GitCommit + ")"
}
